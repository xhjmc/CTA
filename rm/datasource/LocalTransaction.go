package datasource

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/XH-JMC/cta/common/logs"
	"github.com/XH-JMC/cta/common/sqlparser"
	"github.com/XH-JMC/cta/common/sqlparser/model"
	"github.com/XH-JMC/cta/model/rmmodel"
	"sync"
	"time"
)

type LocalTransaction struct {
	xid           string
	branchId      int64
	resourceId    string
	lockKeys      string
	tx            *sql.Tx
	status        rmmodel.BranchStatus
	sqlParserName string

	undoItems       []*UndoItem
	undoLogStmtOnce sync.Once
}

// When some errors occur during committing, the LocalTransaction will be rollback and this function will return error.
func (ltx *LocalTransaction) Commit() (err error) {
	ltx.status = rmmodel.PhaseOne_Failed
	defer func() {
		ltx.reportBranch()
		if err != nil {
			_ = ltx.tx.Rollback()
		}
	}()

	// insert undo_log before committing branch
	if err = ltx.insertUndoLog(); err != nil {
		err = fmt.Errorf("xid: %s, branchId: %d, resourceId: %s, insert undo_log error: %s", ltx.xid, ltx.branchId, ltx.resourceId, err)
		logs.Info(err)
		return
	}

	// get global lock before committing branch
	err = ltx.globalLock()
	if err != nil {
		err = fmt.Errorf("xid: %s, branchId: %d, resourceId: %s, get global lock error: %s", ltx.xid, ltx.branchId, ltx.resourceId, err)
		logs.Info(err)
		return
	}

	err = ltx.tx.Commit()
	if err != nil {
		err = fmt.Errorf("xid: %s, branchId: %d, commit error: %s", ltx.xid, ltx.branchId, err)
		logs.Info(err)
		return
	}

	ltx.status = rmmodel.PhaseOne_Done
	return
}

func (ltx *LocalTransaction) RollBack() error {
	ltx.status = rmmodel.PhaseOne_Failed
	defer ltx.reportBranch()

	err := ltx.tx.Rollback()
	if err != nil {
		err = fmt.Errorf("xid: %s, branchId: %d, rollback error: %s", ltx.xid, ltx.branchId, err)
		logs.Info(err)
	}
	return err
}

// 根据sql语句，构建SQLParser，并生成对应的stmt和镜像stmt
func (ltx *LocalTransaction) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	prepareStmt, err := ltx.tx.PrepareContext(ctx, query)
	if err != nil {
		err = fmt.Errorf("xid: %s, branchId: %d, PrepareContext error: %s", ltx.xid, ltx.branchId, err)
		logs.Info(err)
		return nil, err
	}

	sqlParser, err := sqlparser.NewSQLParser(ltx.sqlParserName, query)
	if err != nil {
		err = fmt.Errorf("xid: %s, branchId: %d, NewSQLParser error: %s", ltx.xid, ltx.branchId, err)
		logs.Info(err)
	}

	stmt := &Stmt{
		ltx:       ltx,
		stmt:      prepareStmt,
		sqlParser: sqlParser,
	}
	switch sqlParser.GetSQLType() {
	case model.INSERT:
		stmt.beforeImageStmt = nil
		stmt.beforeImageArgsFunc = nil
		stmt.afterImageStmt = nil
		stmt.afterImageArgsFunc = nil
	case model.DELETE:
		parser := sqlParser.(model.SQLDeleteParser)
		imageQuery := fmt.Sprintf("select * from %s where %s", parser.GetTableName(), parser.GetCondition())
		imageStmt, err := ltx.tx.PrepareContext(ctx, imageQuery)
		if err != nil {
			err = fmt.Errorf("xid: %s, branchId: %d, imageQuery: %s, prepare image stmt error: %s", ltx.xid, ltx.branchId, imageQuery, err)
			logs.Info(err)
		}
		imageArgsFunc := func(args []interface{}) []interface{} {
			return args
		}
		stmt.beforeImageStmt = imageStmt
		stmt.beforeImageArgsFunc = imageArgsFunc
		stmt.afterImageStmt = nil
		stmt.afterImageArgsFunc = nil
	case model.UPDATE:
		parser := sqlParser.(model.SQLUpdateParser)
		imageQuery := fmt.Sprintf("select * from %s where %s", parser.GetTableName(), parser.GetCondition())
		imageStmt, err := ltx.tx.PrepareContext(ctx, imageQuery)
		if err != nil {
			err = fmt.Errorf("xid: %s, branchId: %d, imageQuery: %s, prepare image stmt error: %s", ltx.xid, ltx.branchId, imageQuery, err)
			logs.Info(err)
		}
		imageArgsFunc := func(args []interface{}) []interface{} {
			return args[len(args)-parser.CountPlaceholderInCondition():]
		}
		stmt.beforeImageStmt = imageStmt
		stmt.beforeImageArgsFunc = imageArgsFunc
		stmt.afterImageStmt = imageStmt
		stmt.afterImageArgsFunc = imageArgsFunc
	case model.SELECT:
		stmt.beforeImageStmt = nil
		stmt.beforeImageArgsFunc = nil
		stmt.afterImageStmt = nil
		stmt.afterImageArgsFunc = nil
	}
	return stmt, nil
}

func (ltx *LocalTransaction) Prepare(query string) (*Stmt, error) {
	return ltx.PrepareContext(context.Background(), query)
}

// 复用PrepareContext和Stmt.ExecContext的逻辑
func (ltx *LocalTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := ltx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.ExecContext(ctx, args...)
}

func (ltx *LocalTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return ltx.ExecContext(context.Background(), query, args...)
}

// 复用PrepareContext和Stmt.QueryContext的逻辑
func (ltx *LocalTransaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := ltx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.QueryContext(ctx, args...)
}

func (ltx *LocalTransaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return ltx.QueryContext(context.Background(), query, args...)
}

// 复用PrepareContext和Stmt.QueryRowContext的逻辑
func (ltx *LocalTransaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	stmt, err := ltx.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}
	defer stmt.Close()
	return stmt.QueryRowContext(ctx, args...)
}

func (ltx *LocalTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return ltx.QueryRowContext(context.Background(), query, args...)
}

func (ltx *LocalTransaction) reportBranch() bool {
	err := GetDataSourceManager().BranchReport(context.Background(), rmmodel.AT, ltx.xid, ltx.branchId, ltx.status)
	if err != nil {
		logs.Infof("xid: %s, branchId: %d, report branch error: %s", ltx.xid, ltx.branchId, err)
		return false
	}
	return true
}

func (ltx *LocalTransaction) addUndoItem(undoItem *UndoItem) {
	ltx.undoLogStmtOnce.Do(func() {
		if ltx.undoItems == nil {
			ltx.undoItems = make([]*UndoItem, 0)
		}
	})
	ltx.undoItems = append(ltx.undoItems, undoItem)
}

func (ltx *LocalTransaction) insertUndoLog() error {
	log := &UndoLog{
		Xid:             ltx.xid,
		BranchId:        ltx.branchId,
		UndoItems:       ltx.undoItems,
		LogStatus:       NormalStatus,
		CreateTimestamp: time.Now().UnixNano(),
	}
	return log.Insert(ltx.tx)
}

func (ltx *LocalTransaction) addLockKey(lockKey string) {
	ltx.lockKeys += lockKey + ";"
}

func (ltx *LocalTransaction) globalLock() error {
	return GetDataSourceManager().GlobalLock(context.Background(), rmmodel.AT, ltx.xid, ltx.resourceId, ltx.lockKeys)
}
