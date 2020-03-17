package client

import (
	"context"
	"cta/common/addrgetter"
	"cta/common/pool"
	"cta/logs"
	"cta/rpc/thrift/config"
	"errors"
	_ "net/http/pprof"
)

type ThriftClientPool interface {
	Borrow(ctx context.Context) (*ThriftClient, error)
	Return(ctx context.Context, tClient *ThriftClient) error
}

type StandardThriftClientPool struct {
	pool pool.Pool
}

func StandardThriftClientPoolFactory(po pool.Pool, conf *config.ThriftConfig, getter addrgetter.AddrGetter) *StandardThriftClientPool {
	p := &StandardThriftClientPool{}
	p.Init(po, conf, getter)
	return p
}

func (p *StandardThriftClientPool) Init(po pool.Pool, conf *config.ThriftConfig, addrGetter addrgetter.AddrGetter) {
	p.pool = po
	p.pool.SetObjectFactory(&pool.ObjectFactory{
		Create: func(ctx context.Context) (interface{}, error) {
			tClient, err := NewTClientWithAddr(conf, addrGetter.GetAddr())
			if err != nil {
				logs.Warnf("new TClient err: %s", err)
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
