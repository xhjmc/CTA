package tm

import (
	"context"
	"cta/model/tmmodel"
	"cta/tc"
	"sync"
)

type TransactionManager struct {
}

var (
	transactionManager     *TransactionManager
	transactionManagerOnce sync.Once
)

func GetTransactionManager() *TransactionManager {
	transactionManagerOnce.Do(func() {
		transactionManager = &TransactionManager{}
	})
	return transactionManager
}

func (tm *TransactionManager) TransactionBegin(ctx context.Context) (string, error) {
	return tc.GetTransactionCoordinatorClient().TransactionBegin(ctx)
}

func (tm *TransactionManager) TransactionCommit(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tc.GetTransactionCoordinatorClient().TransactionCommit(ctx, xid)
}

func (tm *TransactionManager) TransactionRollback(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tc.GetTransactionCoordinatorClient().TransactionRollback(ctx, xid)
}

func (tm *TransactionManager) GetTransactionStatus(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tc.GetTransactionCoordinatorClient().GetTransactionStatus(ctx, xid)
}
