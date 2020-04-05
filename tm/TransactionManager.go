package tm

import (
	"context"
	"github.com/XH-JMC/cta/model/tmmodel"
	"github.com/XH-JMC/cta/tc/tcclient"
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
	return tcclient.GetTransactionCoordinatorClient().TransactionBegin(ctx)
}

func (tm *TransactionManager) TransactionCommit(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tcclient.GetTransactionCoordinatorClient().TransactionCommit(ctx, xid)
}

func (tm *TransactionManager) TransactionRollback(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tcclient.GetTransactionCoordinatorClient().TransactionRollback(ctx, xid)
}

func (tm *TransactionManager) GetTransactionStatus(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tcclient.GetTransactionCoordinatorClient().GetTransactionStatus(ctx, xid)
}
