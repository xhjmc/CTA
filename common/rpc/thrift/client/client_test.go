package client_test

import (
	"context"
	"cta/common/logs"
	"cta/common/rpc/thrift/client"
	"cta/common/rpc/thrift/gen-go/rmservice"
	objPool "github.com/jolestar/go-commons-pool/v2"
	_ "net/http/pprof"
	"sync"
	"testing"
	"time"
)

const (
	n = 100
)

// ensure import "github.com/jolestar/go-commons-pool/v2"
var _ = objPool.DefaultMaxTotal

var (
	ctx    = context.Background()
	wg     = sync.WaitGroup{}
	errNum = 0
)

func TestTClientWithPool(t *testing.T) {
	//addrList := []string{"localhost:5460"}
	//addrList = append(addrList, "[::1]:5461")

	//poolConf := objPool.NewDefaultPoolConfig()
	//poolConf.MaxTotal = 10
	//poolConf.MaxIdle = 10
	//
	//thriftConf := config.GetDefaultThriftConfig()
	//thriftConf.DialTimeout = time.Millisecond * 400
	//thriftConf.ReadWriteTimeout = time.Millisecond * 1600
	//thriftConf.Timeout = time.Millisecond * 2000
	//// 0		4.5						5						9			9.5
	//// 0 send	0 timeout and 1 send	1 recv 0 and 2 send		2 recv 1	2 result
	//
	//tClient := client.TClientWithPoolFactory(
	//	client.StandardThriftClientPoolFactory(
	//		pool.ObjectPoolFactory2(poolConf),
	//		thriftConf,
	//		constant.TCServiceName,
	//	), thriftConf.Timeout,
	//)

	tClient := client.TClientWithPoolFactory3("localhost:5460")

	//go func() {
	//	wg.Add(1)
	//	_ = http.ListenAndServe("0.0.0.0:8080", nil)
	//	wg.Done()
	//}()

	now := time.Now()
	wg.Add(n)
	for i := 0; i < n; i++ {
		id := i
		go func() {
			defer wg.Done()

			cli := rmservice.NewResourceManagerBaseServiceClient(tClient)

			resp, err := cli.Ping(ctx, "ping")
			logs.Info(time.Since(now), id, "ping()", resp, err)

			if err != nil {
				errNum++
				//logs.Info("handle err:", err)
			}
		}()
	}

	logs.Info("before wait")
	wg.Wait()
	logs.Info("after wait")
	logs.Info("errNum:", errNum)
}
