package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/XH-JMC/cta/common/discovery"
	"github.com/XH-JMC/cta/common/logs"
	"github.com/XH-JMC/cta/common/pool"
	"github.com/XH-JMC/cta/common/rpc/thrift/config"
)

type ThriftClientPool interface {
	Borrow(ctx context.Context) (*ThriftClient, error)
	Return(ctx context.Context, tClient *ThriftClient) error
}

type StandardThriftClientPool struct {
	pool pool.Pool
}

func StandardThriftClientPoolFactory(po pool.Pool, conf *config.ThriftConfig,
	serviceName string) *StandardThriftClientPool {
	p := &StandardThriftClientPool{}
	p.Init(po, conf, serviceName)
	return p
}

func (p *StandardThriftClientPool) Init(po pool.Pool, conf *config.ThriftConfig, serviceName string) {
	p.pool = po
	p.pool.SetObjectFactory(&pool.ObjectFactory{
		Create: func(ctx context.Context) (interface{}, error) {
			addr, err := discovery.GetAddr(serviceName)
			if err != nil {
				err = fmt.Errorf("service discovery error: %s", err)
				logs.Warn(err)
				return nil, err
			}
			tClient, err := NewTClientWithAddr(conf, addr)
			if err != nil {
				err = fmt.Errorf("new TClient err: %s", err)
				logs.Warn(err)
				return nil, err
			}
			return tClient, nil
		},
		Destroy: func(ctx context.Context, obj interface{}) error {
			if tClient, ok := obj.(*ThriftClient); ok {
				return tClient.TTransport.Close()
			}
			return nil
		},
		Validate: func(ctx context.Context, obj interface{}) bool {
			tClient, ok := obj.(*ThriftClient)
			return ok && tClient != nil && tClient.TTransport != nil && tClient.TTransport.IsOpen()
		},
	})
	p.pool.Init(context.Background())
}

func (p *StandardThriftClientPool) Borrow(ctx context.Context) (*ThriftClient, error) {
	obj, err := p.pool.Borrow(ctx)
	if err != nil {
		return nil, err
	}
	tClient, ok := obj.(*ThriftClient)
	if ok {
		return tClient, nil
	}
	return nil, errors.New("object is not *ThriftClient")
}

func (p *StandardThriftClientPool) Return(ctx context.Context, tClient *ThriftClient) error {
	return p.pool.Return(ctx, tClient)
}
