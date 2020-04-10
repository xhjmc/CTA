package tcclient

import (
	"context"
	"github.com/XH-JMC/cta/model/rmmodel"
	"github.com/XH-JMC/cta/model/tmmodel"
)

// only for test
type TCMockClient struct {
}

func (m *TCMockClient) BranchRegister(ctx context.Context, branchType rmmodel.BranchType, xid string, resourceId string, applicationName string) (int64, error) {
	return 10086, nil
}

func (m *TCMockClient) BranchReport(ctx context.Context, branchType rmmodel.BranchType, xid string, branchId int64, status rmmodel.BranchStatus) error {
	return nil
}

func (m *TCMockClient) GlobalLock(ctx context.Context, branchType rmmodel.BranchType, xid string, resourceId string, lockKeys string) error {
	return nil
}

func (m *TCMockClient) TransactionBegin(ctx context.Context) (string, error) {
	return "127.0.0.1:5460:1", nil
}

func (m *TCMockClient) TransactionCommit(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tmmodel.CommitDone, nil
}

func (m *TCMockClient) TransactionRollback(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tmmodel.RollbackDone, nil
}

func (m *TCMockClient) GetTransactionStatus(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tmmodel.UnknownTransactionStatus, nil
}
