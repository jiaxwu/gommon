package pool

import "sync"

// 创建新对象
type NewFunc[T any] func() T

// 清理对象
type ClearFunc[T any] func(T) T

type Pool[T any] struct {
	p         sync.Pool
	clearFunc ClearFunc[T]
}

func New[T any](newFunc NewFunc[T], clearFunc ClearFunc[T]) *Pool[T] {
	if newFunc == nil {
		panic("must be provide NewFunc")
	}
	p := &Pool[T]{
		clearFunc: clearFunc,
	}
	p.p.New = func() any {
		return newFunc()
	}
	return p
}

// 获取对象
func (p *Pool[T]) Get() T {
	return p.p.Get().(T)
}

// 归还对象
func (p *Pool[T]) Put(t T) {
	if p.clearFunc != nil {
		t = p.clearFunc(t)
	}
	p.p.Put(t)
}
