package datasource

import (
	"context"
	"cta/common/eventbus"
	"cta/common/logs"
	"cta/common/publicwaitgroup"
	"cta/common/sqlparser/model"
	"cta/model/eventmodel"
	"cta/model/rmmodel"
	"cta/tc/tcclient"
	"cta/util"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataSourceManager struct {
	dataSourceMap map[string]*DataSource
	lock          sync.RWMutex
}

var (
	dataSourceManager *DataSourceManager
	once              sync.Once
)

func GetDataSourceManager() *DataSourceManager {
	once.Do(func() {
		dataSourceManager = &DataSourceManager{
			dataSourceMap: make(map[string]*DataSource),
		}
		dataSourceManager.init()
	})
	return dataSourceManager
}

func (m *DataSourceManager) init() {
	// 订阅事件
	eventbus.Subscribe(eventmodel.BranchCommit_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.RMInboundEvent); ok && e.BranchType == rmmodel.AT {
			e.BranchStatus, e.Error = m.BranchCommit(ctx, e.BranchType, e.Xid, e.BranchId, e.ResourceId)
		}
	})
	eventbus.Subscribe(eventmodel.BranchRollback_EventName, func(ctx context.Context, event eventbus.Event) {
		if e, ok := event.(*eventmodel.RMInboundEvent); ok && e.BranchType == rmmodel.AT {
			e.BranchStatus, e.Error = m.BranchRollback(ctx, e.BranchType, e.Xid, e.BranchId, e.ResourceId)
		}
	})
}

func (m *DataSourceManager) BranchCommit(ctx context.Context, branchType rmmodel.BranchType, xid string, branchId int64, resourceId string) (rmmodel.BranchStatus, error) {
	// 异步删除undo_log
	publicwaitgroup.Add(1)
	go func() {
		defer publicwaitgroup.Done()
		tx, err := m.MustGetDataSource(resourceId).getSQLDB().BeginTx(ctx, nil)
		if err == nil {
			log := &UndoLog{
				Xid:      xid,
				BranchId: branchId,
			}
			err = log.Delete(tx)
			if err != nil {
				_ = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}
		if err != nil {
			logs.Infof("delete undo_log asynchronously error: %s", err)
		}
	}()
	return rmmodel.PhaseTwo_CommitDone, nil
}

func (m *DataSourceManager) BranchRollback(ctx context.Context, branchType rmmodel.BranchType, xid string, branchId int64, resourceId string) (rmmodel.BranchStatus, error) {
	tx, err := m.MustGetDataSource(resourceId).getSQLDB().BeginTx(ctx, nil)
	if err != nil {
		return rmmodel.PhaseTwo_RollbackFailed_Retryable, err
	}

	undoLog := &UndoLog{
		Xid:      xid,
		BranchId: branchId,
	}
	err = undoLog.Select(tx)
	if err != nil {
		// undo_log中没有该日志，证明对应本地分支事务未提交。
		// 为防止资源悬挂，向undo_log表中插入一条log_status=RollbackDoneStatus的日志。
		undoLog := &UndoLog{
			Xid:             xid,
			BranchId:        branchId,
			UndoItems:       nil,
			LogStatus:       RollbackDoneStatus,
			CreateTimestamp: time.Now().UnixNano(),
		}
		err = undoLog.Insert(tx)
		if err != nil {
			_ = tx.Rollback()
			logs.Infof("insert log_status=RollbackDoneStatus into undo_log error: %s", err)
			return rmmodel.PhaseTwo_RollbackFailed_Retryable, err
		} else {
			_ = tx.Commit()
			return rmmodel.PhaseTwo_RollbackDone, nil
		}
	}

	// 建立保存点Before_Undo
	if undoLog.LogStatus != RollbackFailedStatus {
		_, _ = tx.Exec("SAVEPOINT Before_Undo")
	}

	err = m.rollbackUndoLog(ctx, tx, undoLog)
	if err == nil {
		// 根据id删除undo_log表中的对应行
		err = undoLog.DeleteById(tx)
	}
	if err != nil {
		if undoLog.LogStatus != RollbackFailedStatus {
			// 回滚至保存点Before_Undo
			_, _ = tx.Exec("ROLLBACK TO SAVEPOINT Before_Undo")
			//将log_status更改为RollbackFailedStatus
			undoLog.LogStatus = RollbackFailedStatus
			err := undoLog.UpdateLogStatus(tx)
			if err != nil {
				logs.Infof("update log_status=RollbackFailedStatus error: %s", err)
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		} else {
			_ = tx.Rollback()
		}
		return rmmodel.PhaseTwo_RollbackFailed_Retryable, err
	}
	_ = tx.Commit()
	return rmmodel.PhaseTwo_RollbackDone, nil
}

func (m *DataSourceManager) rollbackUndoLog(ctx context.Context, tx *sql.Tx, undoLog *UndoLog) error {
	for _, undoItem := range undoLog.UndoItems {
		var err error
		switch undoItem.SQLType {
		case model.INSERT:
			err = m.rollbackInsertUndoItem(ctx, tx, undoItem)
		case model.DELETE:
			err = m.rollbackDeleteUndoItem(ctx, tx, undoItem)
		case model.UPDATE:
			err = m.rollbackUpdateUndoItem(ctx, tx, undoItem)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *DataSourceManager) rollbackUpdateUndoItem(ctx context.Context, tx *sql.Tx, undoItem *UndoItem) error {
	err := m.checkUpdateAfterImage(ctx, tx, undoItem)
	if err != nil {
		return err
	}
	return m.rollbackUpdateByBeforeImage(ctx, tx, undoItem)
}

// 根据前镜像回滚数据
func (m *DataSourceManager) rollbackUpdateByBeforeImage(ctx context.Context, tx *sql.Tx, undoItem *UndoItem) error {
	if len(undoItem.AfterImage.Rows) > 0 {
		// 获取所有的列，前提保证所有行中的列都一致
		firstRow := undoItem.BeforeImage.Rows[0]
		colNameList := make([]string, 0)
		set := ""
		prefix := ""
		for colName := range firstRow {
			colNameList = append(colNameList, colName)
			set += fmt.Sprintf("%s%s = ?", prefix, colName)
			prefix = ","

		}
		colLen := len(colNameList)
		query := fmt.Sprintf("update %s set %s where %s = ?", undoItem.TableName, set, BusinessTablePK)
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()
		// 填充执行参数
		for _, imageRow := range undoItem.BeforeImage.Rows {
			args := make([]interface{}, 0, colLen+1)
			for _, colName := range colNameList {
				args = append(args, imageRow[colName].Value)
			}
			args = append(args, imageRow[BusinessTablePK].Value)
			_, err := stmt.ExecContext(ctx, args...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 查询后镜像是否与当前数据一致
func (m *DataSourceManager) checkUpdateAfterImage(ctx context.Context, tx *sql.Tx, undoItem *UndoItem) error {
	query := fmt.Sprintf("select * from %s where %s = ?", undoItem.TableName, BusinessTablePK)
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, imageRow := range undoItem.AfterImage.Rows {
		rows, err := stmt.QueryContext(ctx, imageRow[BusinessTablePK].Value)
		if err != nil {
			return err
		}
		colNames, _ := rows.Columns()
		colLen := len(colNames)
		dest := make([]interface{}, colLen, colLen)
		row := make([]interface{}, colLen, colLen)
		for i := 0; i < colLen; i++ {
			dest[i] = &row[i]
		}
		if rows.Next() {
			err := rows.Scan(dest...)
			if err != nil {
				_ = rows.Close()
				return err
			}
			for i, col := range colNames {
				if !util.InterfaceEqual(row[i], imageRow[col].Value) {
					_ = rows.Close()
					return errors.New("after image is not the same as the current data")
				}
			}
		} else {
			_ = rows.Close()
			return errors.New("after image is not the same as the current data")
		}
		_ = rows.Close()
	}
	return nil
}

func (m *DataSourceManager) rollbackDeleteUndoItem(ctx context.Context, tx *sql.Tx, undoItem *UndoItem) error {
	rowsLen := len(undoItem.BeforeImage.Rows)
	if rowsLen == 0 {
		return nil
	}

	// 获取所有的列，前提保证所有的行中的列都一致
	// 同时生成sql语句
	firstRow := undoItem.BeforeImage.Rows[0]
	colNameList := make([]string, 0)
	colNames := ""
	placeholders := ""
	prefix := "("
	for colName := range firstRow {
		colNameList = append(colNameList, colName)
		colNames += prefix + colName
		placeholders += prefix + "?"
		prefix = ","
	}
	colNames += ")"
	placeholders += ")"
	colLen := len(colNameList)
	queryPrefix := fmt.Sprintf("insert into %s%s values", undoItem.TableName, colNames)

	// 批量执行
	maxBatchSize := BatchSize
	if rowsLen < BatchSize {
		maxBatchSize = rowsLen
	}
	batchArgs := make([]interface{}, 0, maxBatchSize*colLen)

	rowsIndex := 0
	batchExec := func(batchTimes, batchSize int) error {
		if batchTimes == 0 || batchSize == 0 {
			return nil
		}

		batchQuery := queryPrefix + placeholders
		for i := 1; i < batchSize; i++ {
			batchQuery += "," + placeholders
		}

		exec := func() error {
			_, err := tx.ExecContext(ctx, batchQuery, batchArgs...)
			return err
		}
		if batchTimes > 1 {
			// 批次数大于1时，使用prepare预编译sql语句，优化执行速度
			stmt, err := tx.PrepareContext(ctx, batchQuery)
			if err != nil {
				return err
			}
			exec = func() error {
				_, err := stmt.ExecContext(ctx, batchArgs...)
				return err
			}
		}

		// 批量执行batchTimes次
		for i := 0; i < batchTimes; i++ {
			// 重置执行参数
			batchArgs = batchArgs[:0]
			// 填充批量执行参数
			for j := 0; j < batchSize; j++ {
				for _, colName := range colNameList {
					batchArgs = append(batchArgs, undoItem.BeforeImage.Rows[rowsIndex][colName].Value)
				}
				rowsIndex++
			}
			// 执行一批
			if err := exec(); err != nil {
				return err
			}
		}

		return nil
	}

	// 以BatchSize为单位，批量执行
	if err := batchExec(rowsLen/BatchSize, BatchSize); err != nil {
		return err
	}
	return batchExec(1, rowsLen%BatchSize)
}

func (m *DataSourceManager) rollbackInsertUndoItem(ctx context.Context, tx *sql.Tx, undoItem *UndoItem) error {
	query := fmt.Sprintf("delete from %s where %s = ?", undoItem.TableName, BusinessTablePK)
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	for _, row := range undoItem.AfterImage.Rows {
		val := row[BusinessTablePK].Value
		_, err := stmt.ExecContext(ctx, val)
		if err != nil {
			logs.Infof("rollback insert %s=%v error: %s", BusinessTablePK, val, err)
		}
	}
	return nil
}

// 向TC注册分支事务，返回branchId
func (m *DataSourceManager) BranchRegister(ctx context.Context, branchType rmmodel.BranchType, xid string, resourceId string, applicationName string) (int64, error) {
	return tcclient.GetTransactionCoordinatorClient().BranchRegister(ctx, branchType, xid, resourceId, applicationName)
}

// 向TC报告分支事务执行状态
func (m *DataSourceManager) BranchReport(ctx context.Context, branchType rmmodel.BranchType, xid string, branchId int64, status rmmodel.BranchStatus) error {
	return tcclient.GetTransactionCoordinatorClient().BranchReport(ctx, branchType, xid, branchId, status)
}

// 向TC申请全局资源锁
func (m *DataSourceManager) GlobalLock(ctx context.Context, branchType rmmodel.BranchType, xid string, resourceId string, lockKeys string) error {
	return tcclient.GetTransactionCoordinatorClient().GlobalLock(ctx, branchType, xid, resourceId, lockKeys)
}

func (m *DataSourceManager) RegisterResource(resource rmmodel.Resource) error {
	dataSource, ok := resource.(*DataSource)
	if !ok {
		return errors.New("only DataSource can be registered in DataSourceManager")
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.dataSourceMap[resource.GetResourceId()] = dataSource
	return nil
}

func (m *DataSourceManager) UnregisterResource(resource rmmodel.Resource) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	resourceId := resource.GetResourceId()
	if _, ok := m.dataSourceMap[resourceId]; !ok {
		return fmt.Errorf("DataSourceManager has no resource with resourceId %s", resourceId)
	}

	delete(m.dataSourceMap, resourceId)
	return nil
}

func (m *DataSourceManager) GetResource(resourceId string) rmmodel.Resource {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.dataSourceMap[resourceId]
}

func (m *DataSourceManager) GetResources() map[string]rmmodel.Resource {
	m.lock.RLock()
	defer m.lock.RUnlock()
	ret := make(map[string]rmmodel.Resource)
	for key, val := range m.dataSourceMap {
		ret[key] = val
	}
	return ret
}

func (m *DataSourceManager) MustGetDataSource(resourceId string) *DataSource {
	return m.GetResource(resourceId).(*DataSource)
}

func (m *DataSourceManager) GetBranchType() rmmodel.BranchType {
	return rmmodel.AT
}
