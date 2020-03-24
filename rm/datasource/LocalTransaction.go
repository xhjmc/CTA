package datasource

import (
	"context"
	"cta/common/logs"
	"cta/rm"
	"database/sql"
	"fmt"
)

type LocalTransaction interface {
	Commit() error
	RollBack() error
	PrepareContext(ctx context.Context, query string) (*Stmt, error)
	Prepare(query string) (*Stmt, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
}

type LocalTx struct {
	xid        string
	branchId   int64
	resourceId string
	lockKeys   string
	tx         *sql.Tx
	status     rm.BranchStatus
}

// When some errors occur during committing, the LocalTransaction will be rollback and this function will return error.
func (tx *LocalTx) Commit() error {
	tx.status = rm.PhaseOne_Failed
	defer tx.reportBranch()

	// get global lock before committing branch
	err := tx.globalLock()
	if err != nil {
		err = fmt.Errorf("xid: %s, branchId: %d, resourceId: %s, get global lock error: %s", tx.xid, tx.branchId, tx.resourceId, err)
		logs.Info(err)
		return err
	}

	err = tx.tx.Commit()
	if err != nil {
		_ = tx.tx.Rollback()
		err = fmt.Errorf("xid: %s, branchId: %d, commit error: %s", tx.xid, tx.branchId, err)
		logs.Info(err)
		return err
	}

	tx.status = rm.PhaseOne_Done
	return nil

}

func (tx *LocalTx) RollBack() error {
	tx.status = rm.PhaseOne_Failed
	defer tx.reportBranch()

	err := tx.tx.Rollback()
	if err != nil {
		err = fmt.Errorf("xid: %s, branchId: %d, rollback error: %s", tx.xid, tx.branchId, err)
		logs.Info(err)
	}
	return err
}

func (tx *LocalTx) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	//todo
	tx.tx.PrepareContext(ctx,query)
}

func (tx *LocalTx) Prepare(query string) (*Stmt, error) {
	return tx.PrepareContext(context.Background(), query)
}

func (tx *LocalTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	//todo
}

func (tx *LocalTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.ExecContext(context.Background(), query, args...)
}

func (tx *LocalTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	//todo
}

func (tx *LocalTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.QueryContext(context.Background(), query, args...)
}

func (tx *LocalTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	//todo

}

func (tx *LocalTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return tx.QueryRowContext(context.Background(), query, args...)
}

func (tx *LocalTx) globalLock() error {
	return GetDataSourceManager().GlobalLock(context.Background(), rm.AT, tx.xid, tx.resourceId, tx.lockKeys)
}

func (tx *LocalTx) reportBranch() bool {
	// todo
	if !ok {
		logs.Warnf("xid: %s, branchId: %d, report branch failed", tx.xid, tx.branchId)
	}
	return ok
}
