package rmmodel

import "context"

type ResourceManagerOutbound interface {
	BranchRegister(ctx context.Context, branchType BranchType, xid string, resourceId string, applicationName string) (int64, error)
	BranchReport(ctx context.Context, branchType BranchType, xid string, branchId int64, status BranchStatus) error
	GlobalLock(ctx context.Context, branchType BranchType, xid string, resourceId string, lockKeys string) error
}
