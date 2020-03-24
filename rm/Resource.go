package rm

type Resource interface {
	GetResourceGroupId() string
	GetResourceId() string
}
