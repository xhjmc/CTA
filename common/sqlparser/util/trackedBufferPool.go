package util

import (
	"github.com/xwb1989/sqlparser"
	"sync"
)

type TrackedBufferPool struct {
	pool sync.Pool
}

var (
	trackedBufferPool     *TrackedBufferPool
	trackedBufferPoolOnce sync.Once
)

func GetTrackedBufferPool() *TrackedBufferPool {
	trackedBufferPoolOnce.Do(func() {
		trackedBufferPool = &TrackedBufferPool{pool: sync.Pool{
			New: func() interface{} {
				return newTrackedBuffer()
			},
		}}
	})
	return trackedBufferPool
}

func (p *TrackedBufferPool) Get() *sqlparser.TrackedBuffer {
	return p.pool.Get().(*sqlparser.TrackedBuffer)
}

func (p *TrackedBufferPool) Put(buf *sqlparser.TrackedBuffer) {
	if buf != nil {
		buf.Reset()
		p.pool.Put(buf)
	}
}

func (p *TrackedBufferPool) Handle(handler func(*sqlparser.TrackedBuffer)) {
	buf := p.Get()
	handler(buf)
	p.Put(buf)
}

func newTrackedBuffer() *sqlparser.TrackedBuffer {
	return sqlparser.NewTrackedBuffer(func(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
		if node, ok := node.(*sqlparser.SQLVal); ok {
			switch node.Type {
			case sqlparser.ValArg:
				buf.WriteArg("?")
				return
			}
		}
		node.Format(buf)
	})
}
