package client

import (
	"context"
	"fmt"
	"github.com/XH-JMC/cta/common/pool"
	"github.com/XH-JMC/cta/common/rpc/thrift/config"
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

func TClientWithPoolFactory2(conf *config.ThriftConfig, serviceName string) *TClientWithPool {
	if conf == nil {
		conf = config.GetDefaultThriftConfig()
	}
	return TClientWithPoolFactory(
		StandardThriftClientPoolFactory(pool.ObjectPoolFactory3(), conf, serviceName),
		conf.Timeout)
}

func TClientWithPoolFactory3(serviceName string) *TClientWithPool {
	return TClientWithPoolFactory2(nil, serviceName)
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
