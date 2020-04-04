package rmmodel

type ResourceManager interface {
	ResourceManagerInbound
	ResourceManagerOutbound
	RegisterResource(resource Resource) error
	UnregisterResource(resource Resource) error
	GetResource(resourceId string) Resource
	GetResources() map[string]Resource // resourceId -> Resource
	GetBranchType() BranchType
}
