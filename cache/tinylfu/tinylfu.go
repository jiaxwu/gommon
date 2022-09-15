package tinylfu

import (
	"fmt"

	"github.com/jiaxwu/gommon/cache"
	"github.com/jiaxwu/gommon/cache/lru"
	"github.com/jiaxwu/gommon/cache/slru"
	"github.com/jiaxwu/gommon/counter/cm"
	"github.com/jiaxwu/gommon/filter/bloom"
)

const (
	// 窗口缓存比例
	windowPercentage = 0.01
	// 过滤器错误率
	filterFalsePositiveRate = 0.01
	// 计数器错误范围
	counterErrorRange = 1
	// 计数器错误率
	counterErrorRate = 0.01
	// 采样因子
	samplesFactor = 8
)

// 转换成[]byte
type BytesFunc[T comparable] func(t T) []byte

// W-TinyLFU
// 非线程安全，请根据业务加锁
// https://arxiv.org/pdf/1512.00727v2.pdf
type Cache[K comparable, V any] struct {
	filter           *bloom.Filter     // 过滤器
	counter          *cm.Counter4      // 计数器
	window           *lru.Cache[K, V]  // 窗口缓存
	main             *slru.Cache[K, V] // 主缓存
	samplesThreshold uint64            // 采样阈值，到达阈值计数会减半
	samples          uint64            // 当前采样数量
	bytesFunc        BytesFunc[K]      // 把Key转换成Bytes的函数
}

func New[K comparable, V any](bytesFunc BytesFunc[K], capacity int) *Cache[K, V] {
	windowCap := int(0.01 * float64(capacity))
	mainCap := capacity - windowCap

	return &Cache[K, V]{
		filter:           bloom.New(uint64(capacity), filterFalsePositiveRate),
		counter:          cm.New4(uint64(capacity), counterErrorRange, counterErrorRate),
		window:           lru.New[K, V](windowCap),
		main:             slru.New[K, V](mainCap),
		samplesThreshold: uint64(capacity) * samplesFactor,
		bytesFunc:        bytesFunc,
	}
}

// 设置 OnEvict
func (c *Cache[K, V]) SetOnEvict(onEvict cache.OnEvict[K, V]) {
	// 只有淘汰段的元素才会真正被淘汰，保护段的元素会先被淘汰到淘汰段
	c.probation.SetOnEvict(onEvict)
}

// 添加或更新元素
func (c *Cache[K, V]) Put(key K, value V) {
	// 先看是否已经在保护段，如果是则更新即可
	if c.protected.Contains(key) {
		c.protected.Put(key, value)
		return
	}

	// 如果在淘汰段，则移动到保护段
	if c.probation.Contains(key) {
		c.moveToProtected(key, value)
		return
	}

	// 如果已经到达最大尺寸，先剔除淘汰段的一个元素
	if c.Full() {
		c.probation.Evict()
	}

	// 添加元素到淘汰段
	c.probation.Put(key, value)
}

// 获取元素
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// 先看是否已经在保护段，如果是则更新即可
	if value, ok := c.protected.Get(key); ok {
		return value, ok
	}

	// 如果在淘汰段，则移动到保护段
	if value, ok := c.probation.Get(key); ok {
		c.moveToProtected(key, value)
		return value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 获取元素，不更新状态
func (c *Cache[K, V]) Peek(key K) (V, bool) {
	if value, ok := c.protected.Peek(key); ok {
		return value, ok
	}

	if value, ok := c.probation.Peek(key); ok {
		return value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 是否包含元素，不更新状态
func (c *Cache[K, V]) Contains(key K) bool {
	return c.probation.Contains(key) || c.protected.Contains(key)
}

// 获取缓存的Keys
func (c *Cache[K, V]) Keys() []K {
	fmt.Println("probation", c.probation.Keys())
	fmt.Println("protected", c.protected.Keys())
	return append(c.probation.Keys(), c.protected.Keys()...)
}

// 获取缓存的Values
func (c *Cache[K, V]) Values() []V {
	return append(c.probation.Values(), c.protected.Values()...)
}

// 获取缓存的Entries
func (c *Cache[K, V]) Entries() []*cache.Entry[K, V] {
	return append(c.probation.Entries(), c.protected.Entries()...)
}

// 移除元素
func (c *Cache[K, V]) Remove(key K) bool {
	if c.protected.Remove(key) {
		return true
	}
	if c.probation.Remove(key) {
		return true
	}
	return false
}

// 淘汰元素
func (c *Cache[K, V]) Evict() *cache.Entry[K, V] {
	if entry := c.probation.Evict(); entry != nil {
		return entry
	}
	if entry := c.protected.Evict(); entry != nil {
		return entry
	}
	return nil
}

// 清空缓存
func (c *Cache[K, V]) Clear(needOnEvict bool) {
	c.probation.Clear(needOnEvict)
	c.protected.Clear(needOnEvict)
}

// 改变容量
func (c *Cache[K, V]) Resize(capacity int, needOnEvict bool) {
	_, protectedCap := splitCap(capacity)
	c.probation.Resize(capacity, needOnEvict)
	c.protected.Resize(protectedCap, needOnEvict)
}

// 元素个数
func (c *Cache[K, V]) Len() int {
	return c.probation.Len() + c.protected.Len()
}

// 容量
func (c *Cache[K, V]) Cap() int {
	return c.probationCap + c.protectedCap
}

// 缓存满了
func (c *Cache[K, V]) Full() bool {
	return c.Len() == c.Cap()
}

// 移动到保护段
func (c *Cache[K, V]) moveToProtected(key K, value V) {
	// 从淘汰段移除
	c.probation.Remove(key)

	// 如果保护段满了，则把保护段的一个元素移动到淘汰段
	if c.protected.Len() == c.protectedCap {
		// 从保护段淘汰一个元素
		entry := c.protected.Evict()
		// 添加到淘汰段
		c.probation.Put(entry.Key, entry.Value)
	}

	// 添加到保护段
	c.protected.Put(key, value)
}

// 分割容量
func splitCap(capacity int) (int, int) {
	protectedCap := int(float64(capacity) * ProtectedPercentage)
	probationCap := capacity - protectedCap
	return probationCap, protectedCap
}
