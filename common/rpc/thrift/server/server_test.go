package server_test

import (
	"cta/common/logs"
	"cta/common/rpc/thrift/gen-go/rmservice"
	"cta/common/rpc/thrift/handler"
	"cta/common/rpc/thrift/server"
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
			processor := rmservice.NewResourceManagerBaseServiceProcessor(&handler.ResourceManagerBaseServiceHandler{})

			s := server.NewThriftServer(addr, processor)
			s.SetConf(nil)
			s.Init()
			if err := s.Run(); err != nil {
				logs.Warnf("server run err: %s", err)
			}
		}()
	}
	wg.Wait()
}
