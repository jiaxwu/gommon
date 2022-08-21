package lru

import "github.com/jiaxwu/gommon/container/list"

type kv[K comparable, V any] struct {
	k K
	v V
}

type LRU[K comparable, V any] struct {
	l    *list.List[kv[K, V]]
	m    map[K]*list.Element[kv[K, V]]
	size int
}

func New[K comparable, V any](size int) *LRU[K, V] {
	return &LRU[K, V]{
		l:    list.New[kv[K, V]](),
		m:    make(map[K]*list.Element[kv[K, V]], size),
		size: size,
	}
}

// Put 添加或更新元素
func (l *LRU[K, V]) Put(k K, v V) {
	// 如果k已经存在，直接把它移到最后面，然后设置新值
	if elem, ok := l.m[k]; ok {
		l.l.MoveToBack(elem)
		elem.Value.v = v
		return
	}

	// 如果已经到达最大尺寸，先剔除一个元素
	if l.l.Len() == l.size {
		front := l.l.Front()
		l.l.Remove(front)
		delete(l.m, front.Value.k)
	}

	// 添加元素
	elem := l.l.PushBack(kv[K, V]{
		k: k,
		v: v,
	})
	l.m[k] = elem
}

// Get 获取元素
func (l *LRU[K, V]) Get(k K) (V, bool) {
	// 如果存在移动到尾部，然后返回
	if elem, ok := l.m[k]; ok {
		l.l.MoveToBack(elem)
		return elem.Value.v, true
	}

	// 不存在返回空值和false
	var v V
	return v, false
}
