package nearlylru

import (
	"time"

	"github.com/jiaxwu/gommon/cache"
)

// 最小采样个数
const MinSamples = 5

type lastAccessEntry[K comparable, V any] struct {
	entry      *cache.Entry[K, V]
	lastAccess time.Time // 最后一次使用时间
}

// 近似最近最少使用
// 基于随机采样
// 优点：不需要额外链表
// 非线程安全，请根据业务加锁
type Cache[K comparable, V any] struct {
	entries  map[K]*lastAccessEntry[K, V]
	capacity int
	samples  int
	onEvict  cache.OnEvict[K, V]
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity < MinSamples {
		panic("too small capacity")
	}
	return &Cache[K, V]{
		entries:  make(map[K]*lastAccessEntry[K, V]),
		capacity: capacity,
		samples:  5,
	}
}

// 设置 OnEvict
func (c *Cache[K, V]) SetOnEvict(onEvict cache.OnEvict[K, V]) {
	c.onEvict = onEvict
}

// 设置采样个数
func (c *Cache[K, V]) SetSamples(samples int) {
	if samples < MinSamples {
		samples = MinSamples
	}
	if c.Cap() < samples {
		panic("too large samples")
	}
	c.samples = samples
}

// 添加或更新元素
// 返回被淘汰的元素
func (c *Cache[K, V]) Put(key K, value V) *cache.Entry[K, V] {
	// 如果 key 已经存在，直接设置新值
	if entry, ok := c.entries[key]; ok {
		entry.entry.Value = value
		entry.lastAccess = time.Now()
		return nil
	}

	// 如果已经到达最大尺寸，先剔除一个元素
	var evicted *cache.Entry[K, V]
	if c.Full() {
		evicted = c.Evict()
	}

	// 添加元素
	c.entries[key] = &lastAccessEntry[K, V]{
		entry: &cache.Entry[K, V]{
			Key:   key,
			Value: value,
		},
		lastAccess: time.Now(),
	}
	return evicted
}

// 获取元素
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// 如果存在更新时间，然后返回
	if entry, ok := c.entries[key]; ok {
		entry.lastAccess = time.Now()
		return entry.entry.Value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 获取元素，不更新状态
func (c *Cache[K, V]) Peek(key K) (V, bool) {
	// 如果存在
	if entry, ok := c.entries[key]; ok {
		return entry.entry.Value, true
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
	for _, entry := range c.entries {
		values[i] = entry.entry.Value
		i++
	}
	return values
}

// 获取缓存的Entries
func (c *Cache[K, V]) Entries() []*cache.Entry[K, V] {
	entries := make([]*cache.Entry[K, V], c.Len())
	i := 0
	for _, entry := range c.entries {
		entries[i] = entry.entry
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
func (c *Cache[K, V]) Evict() *cache.Entry[K, V] {
	var evictEntry *lastAccessEntry[K, V]
	i := 0
	for _, entry := range c.entries {
		if i >= c.samples {
			break
		}
		if evictEntry == nil || entry.lastAccess.Before(evictEntry.lastAccess) {
			evictEntry = entry
		}
		i++
	}
	if evictEntry == nil {
		return nil
	}
	delete(c.entries, evictEntry.entry.Key)
	// 回调
	if c.onEvict != nil {
		c.onEvict(evictEntry.entry)
	}
	return evictEntry.entry
}

// 清空缓存
func (c *Cache[K, V]) Clear(needOnEvict bool) {
	// 触发回调
	if needOnEvict && c.onEvict != nil {
		for _, entry := range c.entries {
			c.onEvict(entry.entry)
		}
	}

	// 清空
	c.entries = make(map[K]*lastAccessEntry[K, V])
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
