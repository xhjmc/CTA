package discovery

type Discovery interface {
	GetAddr(name string) (string, error)
	MustGetAddr(name string) string
	GetAddrList(name string) ([]string, error)
	MustGetAddrList(name string) []string
}

var serviceDiscovery Discovery = DefaultDiscovery()

func SetDiscovery(discovery Discovery) {
	serviceDiscovery = discovery
}

func GetAddr(name string) (string, error) {
	return serviceDiscovery.GetAddr(name)
}

func MustGetAddr(name string) string {
	return serviceDiscovery.MustGetAddr(name)
}

func GetAddrList(name string) ([]string, error) {
	return serviceDiscovery.GetAddrList(name)
}

func MustGetAddrList(name string) []string {
	return serviceDiscovery.MustGetAddrList(name)
}
