package pool

import (
	"context"
	"errors"
	"github.com/jolestar/go-commons-pool/v2"
	"sync"
)

type ObjectFactory struct {
	Create   func(context.Context) (interface{}, error)
	Destroy  func(ctx context.Context, obj interface{}) error
	Validate func(ctx context.Context, obj interface{}) bool
}

type Pool interface {
	SetObjectFactory(objFactory *ObjectFactory)
	Init(ctx context.Context)
	Borrow(ctx context.Context) (interface{}, error)
	Return(ctx context.Context, obj interface{}) error
	Clear(ctx context.Context)
	Close(ctx context.Context)
	IsClosed() bool
}

func Handle(p Pool, handle func(interface{})) error {
	return HandleWithCtx(context.Background(), p, handle)
}

func HandleWithCtx(ctx context.Context, p Pool, handle func(interface{})) error {
	item, err := p.Borrow(ctx)
	if err == nil && item != nil {
		handle(item)
		err = p.Return(ctx, item)
	}
	return err
}

type ObjectPool struct {
	objFactory *ObjectFactory
	objectPool *pool.ObjectPool
	conf       *pool.ObjectPoolConfig
	once       sync.Once
}

func ObjectPoolFactory0(objFactory *ObjectFactory, conf *pool.ObjectPoolConfig) *ObjectPool {
	p := &ObjectPool{}
	p.SetObjectFactory(objFactory)
	p.SetObjectPoolConfig(conf)
	return p
}

func ObjectPoolFactory1(objFactory *ObjectFactory) *ObjectPool {
	return ObjectPoolFactory0(objFactory, pool.NewDefaultPoolConfig())
}

func ObjectPoolFactory2(conf *pool.ObjectPoolConfig) *ObjectPool {
	return ObjectPoolFactory0(nil, conf)
}

func ObjectPoolFactory3() *ObjectPool {
	return ObjectPoolFactory0(nil, pool.NewDefaultPoolConfig())
}

func (p *ObjectPool) SetObjectFactory(objFactory *ObjectFactory) {
	p.objFactory = objFactory
}

func (p *ObjectPool) SetObjectPoolConfig(conf *pool.ObjectPoolConfig) {
	p.conf = conf
}

func (p *ObjectPool) Init(ctx context.Context) {
	p.once.Do(func() {
		factory := pool.NewPooledObjectFactory(func(ctx context.Context) (interface{}, error) {
			return p.objFactory.Create(ctx)
		}, func(ctx context.Context, obj *pool.PooledObject) error {
			if p.objFactory.Destroy != nil {
				return p.objFactory.Destroy(ctx, obj.Object)
			}
			return nil
		}, func(ctx context.Context, obj *pool.PooledObject) bool {
			if p.objFactory.Validate != nil {
				return p.objFactory.Validate(ctx, obj.Object)
			}
			return true
		}, func(ctx context.Context, obj *pool.PooledObject) error {
			if p.objFactory.Validate != nil {
				validate := p.objFactory.Validate(ctx, obj.Object)
				if !validate {
					return errors.New("verification failed")
				}
			}
			return nil
		}, func(ctx context.Context, obj *pool.PooledObject) error {
			if p.objFactory.Validate != nil {
				validate := p.objFactory.Validate(ctx, obj.Object)
				if !validate {
					return errors.New("verification failed")
				}
			}
			return nil
		})

		p.objectPool = pool.NewObjectPool(ctx, factory, p.conf)
	})
}

func (p *ObjectPool) Borrow(ctx context.Context) (interface{}, error) {
	p.Init(ctx)
	return p.objectPool.BorrowObject(ctx)
}

func (p *ObjectPool) Return(ctx context.Context, obj interface{}) error {
	p.Init(ctx)
	return p.objectPool.ReturnObject(ctx, obj)
}

func (p *ObjectPool) Clear(ctx context.Context) {
	p.Init(ctx)
	p.objectPool.Clear(ctx)
}

func (p *ObjectPool) Close(ctx context.Context) {
	p.Init(ctx)
	p.objectPool.Close(ctx)
}

func (p *ObjectPool) IsClosed() bool {
	p.Init(context.Background())
	return p.objectPool.IsClosed()
}
