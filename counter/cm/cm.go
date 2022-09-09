package cm

import (
	"math"

	mmath "github.com/jiaxwu/gommon/math"
)

const (
	// 一个计数需要多少位表示
	countBits = 4
	// 一个计数哈希到多少个槽
	counterDepth = 4
	// 采样因子，也就是平均每10次计数reset()一次
	sampleFactor = 10
)

// 种子
var seeds = []uint64{0xc3a5c85c97cb3127, 0xb492b66fbe98f273, 0x9ae16a3b2f90404f, 0xcbf29ce484222325}

// 计数，类似于布隆过滤器，根据哈希映射到多个位置，然后在对应位置进行计数
// 读取时拿对应位置最小的
type Counter struct {
	counters []uint64
	mask     uint64
	// 每计数多少次应该减少计数
	samples uint64
	// Count()的次数
	times uint64
}

func New(width uint64) *Counter {
	width = mmath.RoundUpPowOf2(width) / (64 / counterDepth / countBits)
	if width < 1 {
		width = 1
	}
	return &Counter{
		counters: make([]uint64, width),
		mask:     width - 1,
		samples:  width * sampleFactor,
	}
}

// 根据增加对应位置的计数
func (c *Counter) Count(hash uint64) {
	start := (hash & 3) << 2

	index0 := c.index(hash, 0)
	index1 := c.index(hash, 1)
	index2 := c.index(hash, 2)
	index3 := c.index(hash, 3)

	added := c.count(index0, start)
	added = c.count(index1, start+1) || added
	added = c.count(index2, start+2) || added
	added = c.count(index3, start+3) || added

	if added {
		c.times++
		if c.times == c.samples {
			c.reset()
		}
	}
}

// 估算对应hash的计数
func (c *Counter) Estimate(hash uint64) uint64 {
	start := (hash & 3) << 2
	minCount := uint64(math.MaxUint64)
	for depth := uint64(0); depth < counterDepth; depth++ {
		index := c.index(hash, depth)
		count := (c.counters[index] >> ((start + depth) << 2)) & 0xf
		minCount = mmath.Min(minCount, count)
	}
	return minCount
}

func (c *Counter) count(index, start uint64) bool {
	offset := start << 2
	mask := uint64(countBits*counterDepth-1) << offset
	if c.counters[index]&mask != mask {
		c.counters[index] += 1 << offset
		return true
	}
	return false
}

// 返回counters下标
func (c *Counter) index(hash, depth uint64) uint64 {
	hash = (hash + seeds[depth]) * seeds[depth]
	hash += (hash >> 32)
	return hash & c.mask
}

// 计数减半
func (c *Counter) reset() {
	for i, count := range c.counters {
		if count != 0 {
			c.counters[i] = (count >> 1) & 0x7777777777777777
		}
	}
	c.times = 0
}
