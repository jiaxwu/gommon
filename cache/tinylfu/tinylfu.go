package tinylfu

import (
	"hash/fnv"

	"github.com/jiaxwu/gommon/math"

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
type BytesFunc[T comparable] func(key T) []byte

// W-TinyLFU
// 没有充分测试，只做学习参考
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
	windowCap := math.Max(int(windowPercentage*float64(capacity)), 1)
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

// // 设置 OnEvict
func (c *Cache[K, V]) SetOnEvict(onEvict cache.OnEvict[K, V]) {
	c.main.SetOnEvict(onEvict)
}

// 添加或更新元素
// 返回被淘汰的元素
func (c *Cache[K, V]) Put(key K, value V) *cache.Entry[K, V] {
	// 计算元素哈希值
	hash := c.hash(key)

	// 增加元素计数
	c.inc(hash)

	// 先添加到window
	candidate := c.window.Put(key, value)
	if candidate == nil {
		return nil
	}

	// 获取main里面的最可能被淘汰的元素
	victim := c.main.Victim()
	if victim == nil {
		return c.main.Put(candidate.Key, candidate.Value)
	}

	// candidate和victim进行PK
	candidateFreq := c.estimate(c.hash(candidate.Key))
	victimFreq := c.estimate(c.hash(victim.Key))
	// 如果candidate胜利则加入主缓存
	if candidateFreq > victimFreq {
		return c.main.Put(candidate.Key, candidate.Value)
	}

	// 否则就被淘汰了
	return candidate
}

// 获取元素
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// 计算元素哈希值
	hash := c.hash(key)

	// 增加元素计数
	c.inc(hash)

	// 判断元素是否存在window
	if value, ok := c.window.Get(key); ok {
		return value, true
	}

	// 判断元素是否存在main
	if value, ok := c.main.Get(key); ok {
		return value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 获取元素，不更新状态
func (c *Cache[K, V]) Peek(key K) (V, bool) {
	if value, ok := c.window.Peek(key); ok {
		return value, true
	}

	if value, ok := c.main.Peek(key); ok {
		return value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 是否包含元素，不更新状态
func (c *Cache[K, V]) Contains(key K) bool {
	return c.window.Contains(key) || c.main.Contains(key)
}

// 获取缓存的Keys
func (c *Cache[K, V]) Keys() []K {
	return append(c.window.Keys(), c.main.Keys()...)
}

// 获取缓存的Values
func (c *Cache[K, V]) Values() []V {
	return append(c.window.Values(), c.main.Values()...)
}

// 获取缓存的Entries
func (c *Cache[K, V]) Entries() []*cache.Entry[K, V] {
	return append(c.window.Entries(), c.main.Entries()...)
}

// 移除元素
func (c *Cache[K, V]) Remove(key K) bool {
	removed := false
	removed = c.window.Remove(key) || removed
	removed = c.main.Remove(key) || removed
	return removed
}

// 清空缓存
func (c *Cache[K, V]) Clear(needOnEvict bool) {
	c.window.Clear(needOnEvict)
	c.main.Clear(needOnEvict)
	c.filter.Clear()
	c.counter.Attenuation(0)
	c.samples = 0
}

// 元素个数
func (c *Cache[K, V]) Len() int {
	return c.window.Len() + c.main.Len()
}

// 容量
func (c *Cache[K, V]) Cap() int {
	return c.window.Cap() + c.main.Cap()
}

// 缓存满了
func (c *Cache[K, V]) Full() bool {
	return c.Len() == c.Cap()
}

// 增加元素计数
func (c *Cache[K, V]) inc(hash uint64) {
	c.samples++
	if c.samples == c.samplesThreshold {
		c.filter.Clear()
		c.counter.Attenuation(2)
		c.samples = 0
	}
	if !c.filter.Contains(hash) {
		c.filter.Add(hash)
	} else {
		c.counter.Add(hash, 1)
	}
}

// 估算元素计数
func (c *Cache[K, V]) estimate(hash uint64) uint8 {
	freq := c.counter.Estimate(hash)
	if c.filter.Contains(hash) {
		freq++
	}
	return freq
}

// 计算哈希值
func (c *Cache[K, V]) hash(key K) uint64 {
	keyBytes := c.bytesFunc(key)
	f := fnv.New64()
	f.Write(keyBytes)
	return f.Sum64()
}
