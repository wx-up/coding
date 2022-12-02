package pool

import "sync"

type Pool[T any] struct {
	pool sync.Pool
}

func NewPool[T any](create func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return create()
			},
		},
	}
}

// Get 获取一个元素
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put 放回一个元素
func (p *Pool[T]) Put(v T) {
	p.pool.Put(v)
}
