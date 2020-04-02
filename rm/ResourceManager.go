package rm

type ResourceManager interface {
	ResourceManagerInbound
	ResourceManagerOutbound
	RegisterResource(resource Resource) error
	UnregisterResource(resource Resource) error
	GetResource(resourceId string) Resource
	GetResources() map[string]Resource // resourceId -> Resource
	GetBranchType() BranchType
}

//type AbstractResourceManager struct {
//	lock        sync.Mutex
//	resourceMap map[string]Resource
//}
//
//func (rm *AbstractResourceManager) RegisterResource(resource Resource) {
//	rm.lock.Lock()
//	defer rm.lock.Unlock()
//	rm.resourceMap[resource.GetResourceId()] = resource
//}
//
//func (rm *AbstractResourceManager) UnregisterResource(resource Resource) {
//	rm.lock.Lock()
//	defer rm.lock.Unlock()
//	delete(rm.resourceMap, resource.GetResourceId())
//}
//
//func (rm *AbstractResourceManager) GetResources() map[string]Resource {
//	rm.lock.Lock()
//	defer rm.lock.Unlock()
//	ret := make(map[string]Resource)
//	for key, val := range rm.resourceMap {
//		ret[key] = val
//	}
//	return ret
//}
//
//func (rm *AbstractResourceManager) BeginTx(ctx context.Context, xid, resourceId string) (*tx.LocalTx, error) {
//	x := sql.DB{}
//	tx:=x.Begin()
//	tx.Stmt()
//	db := rm.resourceMap[resourceId]
//	ltx, err := db.BeginTx(ctx, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	branchId, err := RegisterBranch(xid, resourceId)
//	if err != nil {
//		return nil, err
//	}
//
//	tx := &tx.LocalTx{
//	}
//	return tx, nil
//}
