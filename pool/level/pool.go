package level

import (
	"sync"
)

// 创建新对象
type NewFunc[T any] func(l int) T

// 清理对象
type ClearFunc[T any] func(T) T

type LevelFunc[T any] func(t T) int

// 分级的对象池
type Pool[T any] struct {
	pools     []sync.Pool
	levelFunc LevelFunc[T]
	clearFunc ClearFunc[T]
}

func New[T any](newFunc NewFunc[T], clearFunc ClearFunc[T], levelFunc LevelFunc[T], maxLevel int) *Pool[T] {
	if newFunc == nil {
		panic("must be provide NewFunc")
	}
	p := &Pool[T]{
		clearFunc: clearFunc,
		levelFunc: levelFunc,
		pools:     make([]sync.Pool, maxLevel+1),
	}
	for i := 0; i <= maxLevel; i++ {
		i0 := i
		p.pools[i0].New = func() any {
			return newFunc(i0)
		}
	}
	return p
}

// 获取对象
func (p *Pool[T]) Get(l int) T {
	return p.pools[l].Get().(T)
}

// 归还对象
func (p *Pool[T]) Put(b T) {
	if p.clearFunc != nil {
		b = p.clearFunc(b)
	}
	l := p.levelFunc(b)
	p.pools[l].Put(b)
}
