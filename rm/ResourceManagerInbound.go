package rm

import "context"

type ResourceManagerInbound interface {
	BranchCommit(ctx context.Context, branchType BranchType, xid string, branchId int64, resourceId string) (BranchStatus, error)
	BranchRollback(ctx context.Context, branchType BranchType, xid string, branchId int64, resourceId string) (BranchStatus, error)
}
