package lru

import "github.com/jiaxwu/gommon/container/list"

// 淘汰时触发
type OnEvict[K comparable, V any] func(key K, value V)

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// 最近最少使用
// 非线程安全，请根据业务加锁
type LRU[K comparable, V any] struct {
	evictList *list.List[*Entry[K, V]]
	entries   map[K]*list.Element[*Entry[K, V]]
	capacity  int
	onEvict   OnEvict[K, V]
}

func New[K comparable, V any](capacity int) *LRU[K, V] {
	return &LRU[K, V]{
		evictList: list.New[*Entry[K, V]](),
		entries:   make(map[K]*list.Element[*Entry[K, V]]),
		capacity:  capacity,
	}
}

// 设置 OnEvict
func (c *LRU[K, V]) SetOnEvict(onEvict OnEvict[K, V]) {
	c.onEvict = onEvict
}

// 添加或更新元素
func (c *LRU[K, V]) Put(key K, value V) {
	// 如果 key 已经存在，直接把它移到最前面，然后设置新值
	if elem, ok := c.entries[key]; ok {
		c.evictList.MoveToFront(elem)
		elem.Value.Value = value
		return
	}

	// 如果已经到达最大尺寸，先剔除一个元素
	if c.evictList.Len() == c.capacity {
		c.RemoveOldest()
	}

	// 添加元素
	elem := c.evictList.PushFront(&Entry[K, V]{
		Key:   key,
		Value: value,
	})
	c.entries[key] = elem
}

// 获取元素
func (c *LRU[K, V]) Get(key K) (V, bool) {
	// 如果存在移动到头部，然后返回
	if elem, ok := c.entries[key]; ok {
		c.evictList.MoveToFront(elem)
		return elem.Value.Value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 获取元素，不更新状态
func (c *LRU[K, V]) Peek(key K) (V, bool) {
	// 如果存在
	if elem, ok := c.entries[key]; ok {
		return elem.Value.Value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 获取最老元素，不更新状态
func (c *LRU[K, V]) GetOldest() (*Entry[K, V], bool) {
	elem := c.evictList.Back()
	if elem != nil {
		return elem.Value, true
	}
	return nil, false
}

// 获取最新元素，不更新状态
func (c *LRU[K, V]) GetNewest() (*Entry[K, V], bool) {
	elem := c.evictList.Front()
	if elem != nil {
		return elem.Value, true
	}
	return nil, false
}

// 是否包含元素，不更新状态
func (c *LRU[K, V]) Contains(key K) bool {
	_, ok := c.entries[key]
	return ok
}

// 获取缓存的Keys
// 从老到新
func (c *LRU[K, V]) Keys() []K {
	keys := make([]K, c.Len())
	for ent, i := c.evictList.Back(), 0; ent != nil; ent, i = ent.Prev(), i+1 {
		keys[i] = ent.Value.Key
	}
	return keys
}

// 获取缓存的Values
// 从老到新
func (c *LRU[K, V]) Values() []V {
	values := make([]V, c.Len())
	for ent, i := c.evictList.Back(), 0; ent != nil; ent, i = ent.Prev(), i+1 {
		values[i] = ent.Value.Value
	}
	return values
}

// 获取缓存的Entries
// 从老到新
func (c *LRU[K, V]) Entries() []*Entry[K, V] {
	entries := make([]*Entry[K, V], c.Len())
	for ent, i := c.evictList.Back(), 0; ent != nil; ent, i = ent.Prev(), i+1 {
		entries[i] = ent.Value
	}
	return entries
}

// 移除元素
func (c *LRU[K, V]) Remove(key K) {
	if elem, ok := c.entries[key]; ok {
		c.removeElement(elem)
	}
}

// 移除最老的元素
func (c *LRU[K, V]) RemoveOldest() {
	elem := c.evictList.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}

// 清空缓存
func (c *LRU[K, V]) Clear(needOnEvict bool) {
	// 触发回调
	if needOnEvict && c.onEvict != nil {
		for _, v := range c.entries {
			entry := v.Value
			c.onEvict(entry.Key, entry.Value)
		}
	}

	// 清空
	c.evictList.Init()
	c.entries = make(map[K]*list.Element[*Entry[K, V]])
}

// 改变容量
func (c *LRU[K, V]) Resize(capacity int, needOnEvict bool) {
	diff := c.Len() - capacity
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.RemoveOldest()
	}
	c.capacity = capacity
}

// 元素个数
func (c *LRU[K, V]) Len() int {
	return c.evictList.Len()
}

// 容量
func (c *LRU[K, V]) Cap() int {
	return c.capacity
}

// 移除给定节点
func (c *LRU[K, V]) removeElement(elem *list.Element[*Entry[K, V]]) {
	c.evictList.Remove(elem)
	entry := elem.Value
	delete(c.entries, entry.Key)
	// 回调
	if c.onEvict != nil {
		c.onEvict(entry.Key, entry.Value)
	}
}