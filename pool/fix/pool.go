package fix

import "github.com/jiaxwu/gommon/pool"

type Pool[T any] struct {
	cache     chan T
	newFunc   pool.NewFunc[T]
	clearFunc pool.ClearFunc[T]
}

// cacheSize: 对象池缓存长度
func NewPool[T any](cacheSize int, newFunc pool.NewFunc[T], clearFunc pool.ClearFunc[T]) *Pool[T] {
	if newFunc == nil {
		panic("must be provide NewFunc")
	}
	if cacheSize < 1 {
		panic("cacheSize cannot less then 1")
	}
	return &Pool[T]{
		cache:     make(chan T, cacheSize),
		newFunc:   newFunc,
		clearFunc: clearFunc,
	}
}

func (p *Pool[T]) Get() T {
	select {
	// 从channel读
	case t := <-p.cache:
		return t
		// 如果channel空则申请一个新的对象
	default:
		return p.newFunc()
	}
}

func (p *Pool[T]) Put(t T) {
	t = p.clearFunc(t)
	select {
	// 放入channel
	case p.cache <- t:
	// channel满了则丢弃
	default:
	}
}
