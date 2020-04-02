package datasource

import (
	"context"
	"cta/common/logs"
	"cta/common/sqlparser/model"
	"cta/rm"
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
	dataSourceManager     *DataSourceManager
	dataSourceManagerOnce sync.Once
)

func GetDataSourceManager() *DataSourceManager {
	dataSourceManagerOnce.Do(func() {
		// init dataSourceManager
		dataSourceManager = &DataSourceManager{
			dataSourceMap: make(map[string]*DataSource),
		}
	})
	return dataSourceManager
}

func (m *DataSourceManager) BranchCommit(ctx context.Context, branchType rm.BranchType, xid string, branchId int64, resourceId string) (rm.BranchStatus, error) {
	// 异步删除undo_log
	go func() {
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
	return rm.PhaseTwo_CommittDone, nil
}

func (m *DataSourceManager) BranchRollback(ctx context.Context, branchType rm.BranchType, xid string, branchId int64, resourceId string) (rm.BranchStatus, error) {
	tx, err := m.MustGetDataSource(resourceId).getSQLDB().BeginTx(ctx, nil)
	if err != nil {
		return rm.PhaseTwo_RollbackFailed_Retryable, err
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
		} else {
			_ = tx.Commit()
		}
		return rm.PhaseTwo_RollbackDone, nil
	}

	// 建立保存点Before_Undo
	_, _ = tx.Exec("SAVEPOINT Before_Undo")

	err = m.rollbackUndoLog(ctx, tx, undoLog)
	if err != nil {
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
		return rm.PhaseTwo_RollbackFailed_Retryable, nil
	}
	_ = tx.Commit()
	return rm.PhaseTwo_RollbackDone, nil
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
		// 获取所有的列，前提保证所有的行中的列都一致
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
		query := fmt.Sprintf("update %s set %s where %s = ?", undoItem.TableName, set, BusinessPK)
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()
		// 填充执行参数
		for _, imageRow := range undoItem.AfterImage.Rows {
			args := make([]interface{}, 0, colLen+1)
			for _, colName := range colNameList {
				args = append(args, imageRow[colName].Value)
			}
			args = append(args, imageRow[BusinessPK].Value)
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
	query := fmt.Sprintf("select * from %s where %s = ?", undoItem.TableName, BusinessPK)
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, imageRow := range undoItem.AfterImage.Rows {
		rows, err := stmt.QueryContext(ctx, imageRow[BusinessPK].Value)
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
	if rowsLen > 0 {
		// 获取所有的列，前提保证所有的行中的列都一致
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

		rowsIndex := 0
		// 以BatchSize为单位，批量执行
		batchQueryTimes := rowsLen / BatchSize
		if batchQueryTimes > 0 {
			batchPlaceholders := placeholders
			for i := 1; i < BatchSize; i++ {
				batchPlaceholders += "," + placeholders
			}
			batchQuery := queryPrefix + batchPlaceholders
			stmt, err := tx.PrepareContext(ctx, batchQuery)
			if err != nil {
				return err
			}
			defer stmt.Close()
			batchArgsLen := BatchSize * colLen
			batchArgs := make([]interface{}, 0, batchArgsLen)
			// 批量执行batchQueryTimes次
			for i := 0; i < batchQueryTimes; i++ {
				// 填充批量执行参数
				batchArgs = batchArgs[:0]
				for j := 0; j < BatchSize; j++ {
					for _, colName := range colNameList {
						batchArgs = append(batchArgs, undoItem.BeforeImage.Rows[rowsIndex][colName].Value)
						rowsIndex++
					}
				}
				_, err := stmt.ExecContext(ctx, batchArgs...)
				if err != nil {
					return err
				}
			}
		}

		// 批量执行剩余行
		batchSize := rowsLen % BatchSize
		if batchSize > 0 {
			batchPlaceholders := placeholders
			for i := 1; i < batchSize; i++ {
				batchPlaceholders += "," + placeholders
			}
			batchQuery := queryPrefix + batchPlaceholders

			batchArgsLen := batchSize * colLen
			batchArgs := make([]interface{}, 0, batchArgsLen)
			for j := 0; j < batchSize; j++ {
				for _, colName := range colNameList {
					batchArgs = append(batchArgs, undoItem.BeforeImage.Rows[rowsIndex][colName].Value)
					rowsIndex++
				}
			}
			_, err := tx.ExecContext(ctx, batchQuery, batchArgs...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *DataSourceManager) rollbackInsertUndoItem(ctx context.Context, tx *sql.Tx, undoItem *UndoItem) error {
	query := fmt.Sprintf("delete from %s where %s = ?", undoItem.TableName, BusinessPK)
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	for _, row := range undoItem.AfterImage.Rows {
		val := row[BusinessPK].Value
		_, err := stmt.ExecContext(ctx, val)
		if err != nil {
			logs.Infof("rollback insert %s=%v error: %s", BusinessPK, val, err)
		}
	}
	return nil
}
func (m *DataSourceManager) BranchRegister(ctx context.Context, branchType rm.BranchType, xid string, resourceId string) (int64, error) {
	// todo
	return 10086, nil
}

func (m *DataSourceManager) BranchReport(ctx context.Context, branchType rm.BranchType, xid string, branchId int64, status rm.BranchStatus) error {
	// todo
	return nil
}git

func (m *DataSourceManager) GlobalLock(ctx context.Context, branchType rm.BranchType, xid string, resourceId string, lockKeys string) error {
	// todo
	return nil
}

func (m *DataSourceManager) RegisterResource(resource rm.Resource) error {
	dataSource, ok := resource.(*DataSource)
	if !ok {
		return errors.New("only DataSource can be registered in DataSourceManager")
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.dataSourceMap[resource.GetResourceId()] = dataSource
	return nil
}

func (m *DataSourceManager) UnregisterResource(resource rm.Resource) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	resourceId := resource.GetResourceId()
	if _, ok := m.dataSourceMap[resourceId]; !ok {
		return fmt.Errorf("DataSourceManager has no resource with resourceId %s", resourceId)
	}

	delete(m.dataSourceMap, resourceId)
	return nil
}

func (m *DataSourceManager) GetResource(resourceId string) rm.Resource {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.dataSourceMap[resourceId]
}

func (m *DataSourceManager) GetResources() map[string]rm.Resource {
	m.lock.RLock()
	defer m.lock.RUnlock()
	ret := make(map[string]rm.Resource)
	for key, val := range m.dataSourceMap {
		ret[key] = val
	}
	return ret
}

func (m *DataSourceManager) MustGetDataSource(resourceId string) *DataSource {
	return m.GetResource(resourceId).(*DataSource)
}

func (m *DataSourceManager) GetBranchType() rm.BranchType {
	return rm.AT
}
