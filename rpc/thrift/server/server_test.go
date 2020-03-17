package server

import (
	"cta/logs"
	"cta/rpc/thrift/gen-go/tc"
	"cta/rpc/thrift/handler"
	"sync"
	"testing"
)

func TestNewServer(t *testing.T) {
	addrList := []string{"0.0.0.0:5460"}
	//addrList = append(addrList, "[::1]:5461")

	wg := sync.WaitGroup{}
	for _, str := range addrList {
		addr := str
		wg.Add(1)
		go func() {
			defer wg.Done()
			processor := tc.NewTCServiceProcessor(&handler.TCServiceHandler{})

			s := NewThriftServer(addr, processor)
			s.SetConf(nil)
			s.Init()
			if err := s.Run(); err != nil {
				logs.Warnf("server run err: %s", err)
			}
		}()
	}
	wg.Wait()
}
