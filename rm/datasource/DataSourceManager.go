package datasource

import (
	"context"
	"cta/rm"
	"errors"
	"fmt"
	"sync"
)

type DataSourceManager struct {
	dataSourceMap map[string]*DataSource
	lock          sync.Mutex
}

var (
	dataSourceManager     *DataSourceManager
	dataSourceManagerOnce sync.Once
)

func GetDataSourceManager() *DataSourceManager {
	dataSourceManagerOnce.Do(func() {
		// init dataSourceManagerOnce
		dataSourceManager = &DataSourceManager{}
	})
	return dataSourceManager
}

func (m *DataSourceManager) BranchCommit(ctx context.Context, branchType rm.BranchType, xid string, branchId int64, resourceId string) (rm.BranchStatus, error) {
	//todo
	return rm.PhaseTwo_CommittDone, nil
}

func (m *DataSourceManager) BranchRollback(ctx context.Context, branchType rm.BranchType, xid string, branchId int64, resourceId string) (rm.BranchStatus, error) {
	//todo
	return rm.PhaseTwo_RollbackDone, nil
}

func (m *DataSourceManager) BranchRegister(ctx context.Context, branchType rm.BranchType, xid string, resourceId string) (int64, error) {
	//todo
	return 0, nil
}

func (m *DataSourceManager) BranchReport(ctx context.Context, branchType rm.BranchType, xid string, branchId int64, status rm.BranchStatus) error {
	//todo
	return nil
}

func (m *DataSourceManager) GlobalLock(ctx context.Context, branchType rm.BranchType, xid string, resourceId string, lockKeys string) error {
	//todo
	return nil
}

func (m *DataSourceManager) RegisterResource(resource rm.Resource) error {
	dataSource, ok := resource.(*DataSource)
	if !ok {
		return errors.New("only DataSource can be registered in DataSourceManager")
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.dataSourceMap[resource.GetResourceId()] = dataSource
	return nil
}

func (m *DataSourceManager) UnregisterResource(resource rm.Resource) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	resourceId := resource.GetResourceId()
	if _, ok := m.dataSourceMap[resourceId]; !ok {
		return fmt.Errorf("DataSourceManager has no resource with resourceId %s", resourceId)
	}

	delete(m.dataSourceMap, resourceId)
	return nil
}

func (m *DataSourceManager) GetResources() map[string]rm.Resource {
	ret := make(map[string]rm.Resource)
	for key, val := range m.dataSourceMap {
		ret[key] = val
	}
	return ret
}

func (m *DataSourceManager) GetBranchType() rm.BranchType {
	return rm.AT
}
