package lfu

import "github.com/jiaxwu/gommon/container/list"

// 淘汰时触发
type OnEvict[K comparable, V any] func(entry *Entry[K, V])

type Entry[K comparable, V any] struct {
	Key       K
	Value     V
	Frequency uint64 // 频率
}

// 最少使用
// 优点：有效保护长期热门数据，不会因为偶发情况被移除
// 缺点：如果访问模式改变，可能会导致某些很少访问的数据难以被置换出去
// 非线程安全，请根据业务加锁
type Cache[K comparable, V any] struct {
	entries   map[K]*list.Element[*Entry[K, V]]
	evictList *list.List[*Entry[K, V]]
	capacity  int
	onEvict   OnEvict[K, V]
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity < 1 {
		panic("too small capacity")
	}
	return &Cache[K, V]{
		entries:   make(map[K]*list.Element[*Entry[K, V]]),
		evictList: list.New[*Entry[K, V]](),
		capacity:  capacity,
	}
}

// 设置 OnEvict
func (c *Cache[K, V]) SetOnEvict(onEvict OnEvict[K, V]) {
	c.onEvict = onEvict
}

// 添加或更新元素
func (c *Cache[K, V]) Put(key K, value V) {
	// 如果 key 已经存在，直接更新淘汰顺序，然后设置新值
	if elem, ok := c.entries[key]; ok {
		c.updateEvictList(elem)
		elem.Value.Value = value
		return
	}

	// 如果已经到达最大尺寸，先剔除一个元素
	if c.Full() {
		c.Evict()
	}

	// 添加元素
	elem := c.evictList.PushBack(&Entry[K, V]{
		Key:       key,
		Value:     value,
		Frequency: 1,
	})
	c.entries[key] = elem
}

// 获取元素
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// 如果存在频率+1，然后返回
	if elem, ok := c.entries[key]; ok {
		entry := elem.Value
		entry.Frequency++
		c.updateEvictList(elem)
		return entry.Value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 获取元素，不更新状态
func (c *Cache[K, V]) Peek(key K) (V, bool) {
	// 如果存在
	if elem, ok := c.entries[key]; ok {
		return elem.Value.Value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 是否包含元素，不更新状态
func (c *Cache[K, V]) Contains(key K) bool {
	_, ok := c.entries[key]
	return ok
}

// 获取缓存的Keys
func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, c.Len())
	for elem, i := c.evictList.Back(), 0; elem != nil; elem, i = elem.Prev(), i+1 {
		keys[i] = elem.Value.Key
	}
	return keys
}

// 获取缓存的Values
func (c *Cache[K, V]) Values() []V {
	values := make([]V, c.Len())
	for elem, i := c.evictList.Back(), 0; elem != nil; elem, i = elem.Prev(), i+1 {
		values[i] = elem.Value.Value
	}
	return values
}

// 获取缓存的Entries
func (c *Cache[K, V]) Entries() []*Entry[K, V] {
	entries := make([]*Entry[K, V], c.Len())
	for elem, i := c.evictList.Back(), 0; elem != nil; elem, i = elem.Prev(), i+1 {
		entries[i] = elem.Value
	}
	return entries
}

// 移除元素
func (c *Cache[K, V]) Remove(key K) bool {
	if elem, ok := c.entries[key]; ok {
		c.removeElement(elem)
		return true
	}
	return false
}

// 淘汰元素
func (c *Cache[K, V]) Evict() {
	elem := c.evictList.Back()
	if elem != nil {
		c.removeElement(elem)
		// 回调
		if c.onEvict != nil {
			c.onEvict(elem.Value)
		}
	}
}

// 清空缓存
func (c *Cache[K, V]) Clear(needOnEvict bool) {
	// 触发回调
	if needOnEvict && c.onEvict != nil {
		for elem, i := c.evictList.Back(), 0; elem != nil; elem, i = elem.Prev(), i+1 {
			c.onEvict(elem.Value)
		}
	}

	// 清空
	c.entries = make(map[K]*list.Element[*Entry[K, V]])
	c.evictList.Init()
}

// 改变容量
func (c *Cache[K, V]) Resize(capacity int, needOnEvict bool) {
	diff := c.Len() - capacity
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.Evict()
	}
	c.capacity = capacity
}

// 元素个数
func (c *Cache[K, V]) Len() int {
	return len(c.entries)
}

// 容量
func (c *Cache[K, V]) Cap() int {
	return c.capacity
}

// 缓存满了
func (c *Cache[K, V]) Full() bool {
	return c.Len() == c.Cap()
}

// 移除给定节点
func (c *Cache[K, V]) removeElement(elem *list.Element[*Entry[K, V]]) {
	c.evictList.Remove(elem)
	entry := elem.Value
	delete(c.entries, entry.Key)
}

// 更新淘汰顺序
func (c *Cache[K, V]) updateEvictList(elem *list.Element[*Entry[K, V]]) {
	if elem.Prev() != nil && elem.Prev().Value.Frequency < elem.Value.Frequency {
		// 移动到前面元素之前
		c.evictList.MoveBefore(elem, elem.Prev())
	}
}
