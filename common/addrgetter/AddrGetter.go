package addrgetter

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type AddrGetter interface {
	GetAddr() string
}

type AddrFuncGetter struct {
	addrFunc func() string
}

func NewAddrFuncGetter(addrFunc func() string) *AddrFuncGetter {
	g := &AddrFuncGetter{addrFunc: addrFunc}
	return g
}

func (g *AddrFuncGetter) GetAddr() string {
	if g.addrFunc != nil {
		return g.addrFunc()
	}
	return ""
}

type AddrListGetter struct {
	addrList []string
}

func NewAddrListGetter(addrList []string) *AddrListGetter {
	g := &AddrListGetter{addrList: addrList}
	return g
}

// get addr from addrList randomly
func (g *AddrListGetter) GetAddr() string {
	if g.addrList == nil {
		return ""
	}

	id := rand.Int() % len(g.addrList)
	return g.addrList[id]
}
