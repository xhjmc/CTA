package tc

import (
	"context"
	"cta/common/rpc/thrift/client"
	"cta/common/rpc/thrift/gen-go/tcservice"
	"cta/model/rmmodel"
	"cta/model/tmmodel"
	"cta/variable"
	"errors"
	"sync"
)

type TransactionCoordinatorClient interface {
	rmmodel.ResourceManagerOutbound
	tmmodel.TransactionManagerOutbound
}

type StandardTCClient struct {
	client *tcservice.TransactionCoordinatorServiceClient
}

var (
	tcClientOnce sync.Once
	tcClient     TransactionCoordinatorClient
)

func SetTransactionCoordinatorClient(client TransactionCoordinatorClient) {
	tcClientOnce.Do(func() {
		tcClient = client
	})
}

func GetTransactionCoordinatorClient() TransactionCoordinatorClient {
	tcClientOnce.Do(func() {
		tClient := client.TClientWithPoolFactory3(variable.TCServiceName)
		tcClient = &StandardTCClient{
			client: tcservice.NewTransactionCoordinatorServiceClient(tClient),
		}

	})
	return tcClient
}

func (c *StandardTCClient) BranchRegister(ctx context.Context, branchType rmmodel.BranchType, xid string, resourceId string, applicationName string) (int64, error) {
	req := tcservice.NewBranchRegisterRequest()
	req.BranchType = int32(branchType)
	req.Xid = xid
	req.ResourceId = resourceId
	resp, err := c.client.BranchRegister(ctx, req)
	if err != nil {
		return 0, err
	}
	if len(resp.Error) > 0 {
		err = errors.New(resp.Error)
	}
	return resp.BranchId, err
}

func (c *StandardTCClient) BranchReport(ctx context.Context, branchType rmmodel.BranchType, xid string, branchId int64, status rmmodel.BranchStatus) error {
	req := tcservice.NewBranchReportRequest()
	req.BranchType = int32(branchType)
	req.Xid = xid
	req.BranchId = branchId
	req.BranchStatus = int32(status)
	resp, err := c.client.BranchReport(ctx, req)
	if err != nil {
		return err
	}
	if len(resp.Error) > 0 {
		err = errors.New(resp.Error)
	}
	return err
}

func (c *StandardTCClient) GlobalLock(ctx context.Context, branchType rmmodel.BranchType, xid string, resourceId string, lockKeys string) error {
	req := tcservice.NewGlobalLockRequest()
	req.BranchType = int32(branchType)
	req.Xid = xid
	req.ResourceId = resourceId
	req.LockKeys = lockKeys
	resp, err := c.client.GlobalLock(ctx, req)
	if err != nil {
		return err
	}
	if len(resp.Error) > 0 {
		err = errors.New(resp.Error)
	}
	return err
}

func (c *StandardTCClient) TransactionBegin(ctx context.Context) (string, error) {
	req := tcservice.NewTransactionBeginRequest()
	resp, err := c.client.TransactionBegin(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Error) > 0 {
		err = errors.New(resp.Error)
	}
	return resp.Xid, err
}

func (c *StandardTCClient) TransactionCommit(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	req := tcservice.NewTransactionRequest()
	req.Xid = xid
	resp, err := c.client.TransactionCommit(ctx, req)
	if err != nil {
		return tmmodel.Unknown, err
	}
	if len(resp.Error) > 0 {
		err = errors.New(resp.Error)
	}
	return tmmodel.TransactionStatus(resp.TransactionStatus), err
}

func (c *StandardTCClient) TransactionRollback(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	req := tcservice.NewTransactionRequest()
	req.Xid = xid
	resp, err := c.client.TransactionRollback(ctx, req)
	if err != nil {
		return tmmodel.Unknown, err
	}
	if len(resp.Error) > 0 {
		err = errors.New(resp.Error)
	}
	return tmmodel.TransactionStatus(resp.TransactionStatus), err
}

func (c *StandardTCClient) GetTransactionStatus(ctx context.Context, xid string) (tmmodel.TransactionStatus, error) {
	req := tcservice.NewTransactionRequest()
	req.Xid = xid
	resp, err := c.client.GetTransactionStatus(ctx, req)
	if err != nil {
		return tmmodel.Unknown, err
	}
	if len(resp.Error) > 0 {
		err = errors.New(resp.Error)
	}
	return tmmodel.TransactionStatus(resp.TransactionStatus), err
}