package pool

import (
	"context"
	"errors"
	"fmt"
	"github.com/jolestar/go-commons-pool/v2"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	yes := true
	_ = yes

	factory := pool.NewPooledObjectFactorySimple(
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(time.Millisecond * 500)
			if ctx.Value("ok") == true {
				return &yes, nil
			}
			return nil, errors.New("create object fail") // simulate the failure of creating objects
		})

	ctx := context.Background()
	//ctx, f := context.WithTimeout(ctx, time.Second)
	//defer f()
	conf := pool.NewDefaultPoolConfig()
	conf.MaxTotal = 1
	p := pool.NewObjectPool(ctx, factory, conf)

	n := 3
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		id := i
		if id == 2 {
			time.Sleep(time.Millisecond * 1000)
		}
		go func() { // when using goroutines without locking BorrowObject, it will deadlock
			defer wg.Done()
			fmt.Println(id, "before borrow")
			obj, err := p.BorrowObject(context.WithValue(ctx, "ok", id == 2))
			fmt.Println(id, "after borrow")
			if err != nil {
				fmt.Println(id, "borrow err:", err)
			}

			if id == 2 {
				time.Sleep(time.Second*5)
			}

			fmt.Println(id, "before return")
			err = p.ReturnObject(ctx, obj)
			fmt.Println(id, "after return")
			if err != nil {
				fmt.Println(id, "return err:", err)
			}
		}()
	}
	wg.Wait()
}
