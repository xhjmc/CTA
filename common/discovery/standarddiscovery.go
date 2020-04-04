package discovery

import "errors"

type StandardDiscovery struct {
	GetAddrFunc     func(name string) (string, error)
	GetAddrListFunc func(name string) ([]string, error)
}

func DefaultDiscovery() *StandardDiscovery {
	return GenerateStandardDiscovery(func(name string) (string, error) {
		return name, nil
	}, func(name string) ([]string, error) {
		return []string{name}, nil
	})
}

func GenerateStandardDiscovery(GetAddrFunc func(name string) (string, error),
	GetAddrListFunc func(name string) ([]string, error)) *StandardDiscovery {
	return &StandardDiscovery{
		GetAddrFunc:     GetAddrFunc,
		GetAddrListFunc: GetAddrListFunc,
	}
}

func (d *StandardDiscovery) GetAddr(name string) (string, error) {
	if d.GetAddrFunc != nil {
		return d.GetAddrFunc(name)
	}
	return "", errors.New("GetAddrFunc is empty")
}

func (d *StandardDiscovery) MustGetAddr(name string) string {
	addr, _ := d.GetAddr(name)
	return addr
}

func (d *StandardDiscovery) GetAddrList(name string) ([]string, error) {
	if d.GetAddrListFunc != nil {
		return d.GetAddrListFunc(name)
	}
	return nil, errors.New("GetAddrListFunc is empty")
}

func (d *StandardDiscovery) MustGetAddrList(name string) []string {
	addrList, _ := d.GetAddrList(name)
	return addrList
}
