package arc

import (
	"github.com/jiaxwu/gommon/math"

	"github.com/jiaxwu/gommon/cache/lfu"
	"github.com/jiaxwu/gommon/cache/lru"
)

// 淘汰时触发
type OnEvict[K comparable, V any] func(entry *Entry[K, V])

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// 结合LRU和LFU，根据负载动态调整LRU和LFU容量
// 优点：集合LRU和LFU优点
// 缺点：比较耗空间
// 非线程安全，请根据业务加锁
type Cache[K comparable, V any] struct {
	lruCache *lru.Cache[K, V]
	lruEvict *lru.Cache[K, V]
	lfuCache *lfu.Cache[K, V]
	lfuEvict *lfu.Cache[K, V]
	// 表示有多偏向LRU
	preferLRU int
	capacity  int
	onEvict   OnEvict[K, V]
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity < 1 {
		panic("too small capacity")
	}
	c := &Cache[K, V]{
		lruCache: lru.New[K, V](capacity),
		lruEvict: lru.New[K, V](capacity),
		lfuCache: lfu.New[K, V](capacity),
		lfuEvict: lfu.New[K, V](capacity),
		capacity: capacity,
	}
	c.lruCache.SetOnEvict(func(entry *lru.Entry[K, V]) {
		c.lruEvict.Put(entry.Key, entry.Value)
		c.doOnEvict(entry.Key, entry.Value)
	})
	c.lfuCache.SetOnEvict(func(entry *lfu.Entry[K, V]) {
		c.lfuEvict.Put(entry.Key, entry.Value)
		c.doOnEvict(entry.Key, entry.Value)
	})
	return c
}

// 设置 OnEvict
func (c *Cache[K, V]) SetOnEvict(onEvict OnEvict[K, V]) {
	c.onEvict = onEvict
}

// 添加或更新元素
func (c *Cache[K, V]) Put(key K, value V) {
	// 如果存在LRUCache，则移动到LFUCache
	if c.lruCache.Contains(key) {
		c.lruCache.Remove(key)
		c.lfuCache.Put(key, value)
		return
	}

	// 如果存在LFUCache，则更新
	if c.lfuCache.Contains(key) {
		c.lfuCache.Put(key, value)
		return
	}

	// 如果存在LRUEvict，则增加LRUCache的权重
	if c.lruEvict.Contains(key) {
		// 不超过容量，每次最少增加1
		c.preferLRU = math.Min(c.Cap(), c.preferLRU+math.Max(c.lfuEvict.Len()/c.lruEvict.Len(), 1))
		if c.Full() {
			c.evict(false)
		}

		// 移动到LFUCache
		c.lruEvict.Remove(key)
		c.lfuCache.Put(key, value)
		return
	}

	// 如果存在LFUEvict，则减少LRUCache的权重
	if c.lfuEvict.Contains(key) {
		// 不超过容量，每次最少增加1
		c.preferLRU = math.Min(c.Cap(), c.preferLRU-math.Max(c.lruEvict.Len()/c.lfuEvict.Len(), 1))
		if c.Full() {
			c.evict(true)
		}

		// 移动到LFUCache
		c.lfuEvict.Remove(key)
		c.lfuCache.Put(key, value)
		return
	}

	// 如果已经到达最大尺寸，先剔除一个元素
	if c.Full() {
		c.evict(false)
	}

	if c.lruEvict.Len() > c.Cap()-c.preferLRU {
		c.lruEvict.Evict()
	}
	if c.lfuEvict.Len() > c.preferLRU {
		c.lfuEvict.Evict()
	}

	// 添加到LRUCache
	c.lruCache.Put(key, value)
}

// 获取元素
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// 如果存在LRUCache，则移动到LFUCache
	if value, ok := c.lruCache.Peek(key); ok {
		c.lruCache.Remove(key)
		c.lfuCache.Put(key, value)
		return value, true
	}
	// 如果存在LFUCache
	if value, ok := c.lfuCache.Get(key); ok {
		return value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 获取元素，不更新状态
func (c *Cache[K, V]) Peek(key K) (V, bool) {
	// 如果存在
	if value, ok := c.lruCache.Peek(key); ok {
		return value, true
	}
	if value, ok := c.lfuCache.Peek(key); ok {
		return value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 是否包含元素，不更新状态
func (c *Cache[K, V]) Contains(key K) bool {
	return c.lruCache.Contains(key) || c.lfuCache.Contains(key)
}

// 获取缓存的Keys
func (c *Cache[K, V]) Keys() []K {
	return append(c.lruCache.Keys(), c.lfuCache.Keys()...)
}

// 获取缓存的Values
func (c *Cache[K, V]) Values() []V {
	return append(c.lruCache.Values(), c.lfuCache.Values()...)
}

// 获取缓存的Entries
func (c *Cache[K, V]) Entries() []*Entry[K, V] {
	entries := make([]*Entry[K, V], c.Len())
	i := 0
	for _, entry := range c.lruCache.Entries() {
		entries[i] = &Entry[K, V]{
			Key:   entry.Key,
			Value: entry.Value,
		}
		i++
	}
	for _, entry := range c.lfuCache.Entries() {
		entries[i] = &Entry[K, V]{
			Key:   entry.Key,
			Value: entry.Value,
		}
		i++
	}
	return entries
}

// 移除元素
func (c *Cache[K, V]) Remove(key K) bool {
	if c.lruCache.Remove(key) {
		return true
	}
	if c.lfuCache.Remove(key) {
		return true
	}
	if c.lruEvict.Remove(key) {
		return true
	}
	if c.lfuEvict.Remove(key) {
		return true
	}
	return false
}

// 清空缓存
func (c *Cache[K, V]) Clear(needOnEvict bool) {
	c.lruCache.Clear(needOnEvict)
	c.lruEvict.Clear(needOnEvict)
	c.lfuCache.Clear(needOnEvict)
	c.lfuEvict.Clear(needOnEvict)
}

// 元素个数
func (c *Cache[K, V]) Len() int {
	return c.lruCache.Len() + c.lfuCache.Len()
}

// 容量
func (c *Cache[K, V]) Cap() int {
	return c.capacity
}

// 缓存满了
func (c *Cache[K, V]) Full() bool {
	return c.Len() >= c.Cap()
}

// 淘汰元素
// lfuEvictContainsKey: 如果lfuEvict包含key，则先从lruCache淘汰
func (c *Cache[K, V]) evict(lfuEvictContainsKey bool) {
	lruCacheLen := c.lruCache.Len()
	if lruCacheLen > 0 && (lruCacheLen > c.preferLRU || (lruCacheLen == c.preferLRU && lfuEvictContainsKey)) {
		c.lruCache.Evict()
	} else {
		c.lfuCache.Evict()
	}
}

func (c *Cache[K, V]) doOnEvict(key K, value V) {
	if c.onEvict != nil {
		c.onEvict(&Entry[K, V]{
			Key:   key,
			Value: value,
		})
	}
}
