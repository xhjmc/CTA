package client

import (
	"context"
	"cta/common/addrgetter"
	"cta/common/pool"
	"cta/rpc/thrift/config"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

type TClientWithPool struct {
	pool    ThriftClientPool
	timeout time.Duration
}

func TClientWithPoolFactory(p ThriftClientPool, timeout time.Duration) *TClientWithPool {
	return &TClientWithPool{pool: p, timeout: timeout}
}

func DefaultTClientWithPoolFactory(conf *config.ThriftConfig, addrFunc func() string, addrList []string) *TClientWithPool {
	if conf == nil {
		conf = config.GetDefaultThriftConfig()
	}
	var getter addrgetter.AddrGetter
	if addrFunc != nil {
		getter = addrgetter.NewAddrFuncGetter(addrFunc)
	}
	if addrList != nil {
		getter = addrgetter.NewAddrListGetter(addrList)
	}
	return TClientWithPoolFactory(StandardThriftClientPoolFactory(pool.ObjectPoolFactory3(), conf, getter), conf.Timeout)
}

func (c *TClientWithPool) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	if _, ok := ctx.Deadline(); !ok {
		ctx, _ = context.WithTimeout(ctx, c.timeout)
	}
	cli, err := c.pool.Borrow(ctx)
	if err != nil {
		return fmt.Errorf("borrow from client pool error: %s", err)
	}
	defer c.pool.Return(ctx, cli)

	return cli.TClient.Call(ctx, method, args, result)
}
