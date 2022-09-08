package slru

import (
	"fmt"

	"github.com/jiaxwu/gommon/cache"
	"github.com/jiaxwu/gommon/cache/lru"
)

// 保护段比例
const ProtectedPercentage = 0.8

// 分段最近最少使用
// 第一次access是淘汰段，第二次access才进入保护段
// 避免某些很少读取的值把一直读取的值给淘汰了
// 优点：稳定淘汰，避免大量失效
// 非线程安全，请根据业务加锁
type Cache[K comparable, V any] struct {
	probation    *lru.Cache[K, V] // 淘汰段
	protected    *lru.Cache[K, V] // 保护段
	probationCap int
	protectedCap int
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	probationCap, protectedCap := splitCap(capacity)
	return &Cache[K, V]{
		probation:    lru.New[K, V](capacity),
		protected:    lru.New[K, V](protectedCap),
		probationCap: probationCap,
		protectedCap: protectedCap,
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
