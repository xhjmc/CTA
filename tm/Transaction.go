package tm

import (
	"context"
	"github.com/XH-JMC/cta/constant"
	"github.com/XH-JMC/cta/model/tmmodel"
)

type Transaction struct {
	xid    string
	ctx    context.Context
	status tmmodel.TransactionStatus
}

func (tx *Transaction) GetXid() string {
	return tx.xid
}

func (tx *Transaction) GetContext() context.Context {
	return tx.ctx
}

func (tx *Transaction) GetStatus() tmmodel.TransactionStatus {
	return tx.status
}

func (tx *Transaction) Load(ctx context.Context, xid string) error {
	tx.xid = xid
	tx.ctx = context.WithValue(ctx, constant.XidKey, xid)
	return tx.LoadStatus()
}

func (tx *Transaction) LoadStatus() (err error) {
	tx.status, err = GetTransactionManager().GetTransactionStatus(tx.ctx, tx.xid)
	return err
}

func (tx *Transaction) Begin(ctx context.Context) (err error) {
	tx.xid, err = GetTransactionManager().TransactionBegin(ctx)
	tx.ctx = context.WithValue(ctx, constant.XidKey, tx.xid)
	tx.status = tmmodel.Begin
	return
}

func (tx *Transaction) Commit() (err error) {
	tx.status, err = GetTransactionManager().TransactionCommit(tx.ctx, tx.xid)
	return err
}

func (tx *Transaction) Rollback() (err error) {
	tx.status, err = GetTransactionManager().TransactionRollback(tx.ctx, tx.xid)
	return err
}

func (tx *Transaction) Exec(f func(ctx context.Context) error) error {
	return f(tx.ctx)
}
