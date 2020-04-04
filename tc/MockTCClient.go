package tc

import (
	"context"
	"cta/model/rmmodel"
	"cta/model/tmmodel"
)

// only for test
type MockTCClient struct {
}

func (m *MockTCClient) BranchRegister(ctx context.Context, branchType rmmodel.BranchType, xid string, resourceId string, applicationName string) (int64, error) {
	return 10086, nil
}

func (m *MockTCClient) BranchReport(ctx context.Context, branchType rmmodel.BranchType, xid string, branchId int64, status rmmodel.BranchStatus) error {
	return nil
}

func (m *MockTCClient) GlobalLock(ctx context.Context, branchType rmmodel.BranchType, xid string, resourceId string, lockKeys string) error {
	return nil
}

func (m *MockTCClient) TransactionBegin(ctx context.Context) (string, error) {
	return "127.0.0.1:5460:1", nil
}

func (m *MockTCClient) TransactionCommit(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tmmodel.CommitDone, nil
}

func (m *MockTCClient) TransactionRollback(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tmmodel.RollbackDone, nil
}

func (m *MockTCClient) GetTransactionStatus(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	return tmmodel.Unknown, nil
}
