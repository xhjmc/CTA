package client_test

import (
	"context"
	"cta/common/addrgetter"
	"cta/common/logs"
	"cta/common/pool"
	"cta/rpc/thrift/client"
	"cta/rpc/thrift/config"
	"cta/rpc/thrift/gen-go/tc"
	"fmt"
	objPool "github.com/jolestar/go-commons-pool/v2"
	_ "net/http/pprof"
	"sync"
	"testing"
	"time"
)

const (
	n = 10000
)

var (
	ctx    = context.Background()
	wg     = sync.WaitGroup{}
	errNum = 0
)

func TestTClientWithPool(t *testing.T) {
	addrList := []string{"localhost:5460"}
	//addrList = append(addrList, "[::1]:5461")

	poolConf := objPool.NewDefaultPoolConfig()
	poolConf.MaxTotal = 10
	poolConf.MaxIdle = 10

	thriftConf := config.GetDefaultThriftConfig()
	thriftConf.DialTimeout = time.Millisecond * 400
	thriftConf.ReadWriteTimeout = time.Millisecond * 1600
	thriftConf.Timeout = time.Millisecond * 2000
	// 0		4.5						5						9			9.5
	// 0 send	0 timeout and 1 send	1 recv 0 and 2 send		2 recv 1	2 result

	tClient := client.TClientWithPoolFactory(
		client.StandardThriftClientPoolFactory(
			pool.ObjectPoolFactory2(poolConf),
			thriftConf,
			addrgetter.NewAddrListGetter(addrList),
		), thriftConf.Timeout,
	)

	//tClient = DefaultTClientWithPoolFactory(nil, nil, addrList)

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

			cli := tc.NewTCServiceClient(tClient)

			req := tc.NewPingRequest()
			req.Msg = fmt.Sprintf("%d ping", id)
			resp, err := cli.Ping(ctx, req)
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
