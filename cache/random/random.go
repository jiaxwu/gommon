package random

import "github.com/jiaxwu/gommon/cache"

// 随机
// 优点：实现简单
// 非线程安全，请根据业务加锁
type Cache[K comparable, V any] struct {
	entries  map[K]V
	capacity int
	onEvict  cache.OnEvict[K, V]
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity < 1 {
		panic("too small capacity")
	}
	return &Cache[K, V]{
		entries:  make(map[K]V),
		capacity: capacity,
	}
}

// 设置 OnEvict
func (c *Cache[K, V]) SetOnEvict(onEvict cache.OnEvict[K, V]) {
	c.onEvict = onEvict
}

// 添加或更新元素
func (c *Cache[K, V]) Put(key K, value V) {
	// 如果 key 已经存在，直接设置新值
	if _, ok := c.entries[key]; ok {
		c.entries[key] = value
		return
	}

	// 如果已经到达最大尺寸，先剔除一个元素
	if c.Full() {
		c.Evict()
	}

	// 添加元素
	c.entries[key] = value
}

// 获取元素
func (c *Cache[K, V]) Get(key K) (V, bool) {
	return c.Peek(key)
}

// 获取元素
func (c *Cache[K, V]) Peek(key K) (V, bool) {
	value, ok := c.entries[key]
	return value, ok
}

// 是否包含元素
func (c *Cache[K, V]) Contains(key K) bool {
	_, ok := c.entries[key]
	return ok
}

// 获取缓存的Keys
func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, c.Len())
	i := 0
	for key := range c.entries {
		keys[i] = key
		i++
	}
	return keys
}

// 获取缓存的Values
func (c *Cache[K, V]) Values() []V {
	values := make([]V, c.Len())
	i := 0
	for _, value := range c.entries {
		values[i] = value
		i++
	}
	return values
}

// 获取缓存的Entries
func (c *Cache[K, V]) Entries() []*cache.Entry[K, V] {
	entries := make([]*cache.Entry[K, V], c.Len())
	i := 0
	for key, value := range c.entries {
		entries[i] = &cache.Entry[K, V]{
			Key:   key,
			Value: value,
		}
		i++
	}
	return entries
}

// 移除元素
func (c *Cache[K, V]) Remove(key K) bool {
	if _, ok := c.entries[key]; ok {
		delete(c.entries, key)
		return true
	}
	return false
}

// 淘汰元素
func (c *Cache[K, V]) Evict() {
	for key, value := range c.entries {
		delete(c.entries, key)
		// 回调
		if c.onEvict != nil {
			c.onEvict(&cache.Entry[K, V]{
				Key:   key,
				Value: value,
			})
		}
		return
	}
}

// 清空缓存
func (c *Cache[K, V]) Clear(needOnEvict bool) {
	// 触发回调
	if needOnEvict && c.onEvict != nil {
		for key, value := range c.entries {
			c.onEvict(&cache.Entry[K, V]{
				Key:   key,
				Value: value,
			})
		}
	}

	// 清空
	c.entries = make(map[K]V)
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
