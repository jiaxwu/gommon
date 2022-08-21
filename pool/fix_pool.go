package pool

type FixPool[T any] struct {
	cache     chan T
	newFunc   NewFunc[T]
	clearFunc ClearFunc[T]
}

// cacheSize: 对象池缓存长度
func NewFixPool[T any](cacheSize int, newFunc NewFunc[T], clearFunc ClearFunc[T]) *FixPool[T] {
	if newFunc == nil {
		panic("must be provide NewFunc")
	}
	if cacheSize < 1 {
		panic("cacheSize cannot less then 1")
	}
	return &FixPool[T]{
		cache:     make(chan T, cacheSize),
		newFunc:   newFunc,
		clearFunc: clearFunc,
	}
}

func (p *FixPool[T]) Get() T {
	select {
	// 从channel读
	case t := <-p.cache:
		return t
		// 如果channel空则申请一个新的对象
	default:
		return p.newFunc()
	}
}

func (p *FixPool[T]) Put(t T) {
	t = p.clearFunc(t)
	select {
	// 放入channel
	case p.cache <- t:
	// channel满了则丢弃
	default:
	}
}
