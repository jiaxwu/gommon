package cm

import (
	"math"

	mmath "github.com/jiaxwu/gommon/math"
)

const (
	// 一个计数需要多少位表示
	countBits uint64 = 4
	// 一个计数哈希到多少个槽
	counterDepth uint64 = 4
	// 计数掩码
	countMask uint64 = countBits*counterDepth - 1
	// 采样因子，也就是平均多少次计数Reset()一次
	sampleFactor uint64 = 8
)

// 种子
var seeds = []uint64{0xc3a5c85c97cb3127, 0xb492b66fbe98f273, 0x9ae16a3b2f90404f, 0xcbf29ce484222325}

// 计数器，原理类似于布隆过滤器，根据哈希映射到多个位置，然后在对应位置进行计数
// 读取时拿对应位置最小的
type Counter struct {
	counters []uint64
	// counters index掩码
	mask uint64
	// 每计数多少次应该减少计数
	samples uint64
	// Add()的次数
	additions uint64
}

func New(width uint64) *Counter {
	// 2的次方才能用掩码取下标
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
func (c *Counter) Add(hash uint64) {
	added := false
	for depth := uint64(0); depth < counterDepth; depth++ {
		index, offset := c.pos(hash, depth)
		added = c.add(index, offset) || added
	}

	if added {
		c.additions++
		if c.additions == c.samples {
			c.Reset()
		}
	}
}

// 估算对应hash的计数
func (c *Counter) Estimate(hash uint64) uint64 {
	minCount := uint64(math.MaxUint64)
	for depth := uint64(0); depth < counterDepth; depth++ {
		index, offset := c.pos(hash, depth)
		count := (c.counters[index] >> offset) & countMask
		minCount = mmath.Min(minCount, count)
	}
	return minCount
}

// 计数减半
func (c *Counter) Reset() {
	for i, count := range c.counters {
		if count != 0 {
			c.counters[i] = (count >> 1) & 0x7777777777777777
		}
	}
	c.additions = 0
}

// 增加对应下标和偏移的计数
func (c *Counter) add(index, offset uint64) bool {
	mask := countMask << offset
	if c.counters[index]&mask != mask {
		c.counters[index] += 1 << offset
		return true
	}
	return false
}

// 返回hash在counters的位置
// index是数组下标
// offset是对应元素的偏移
func (c *Counter) pos(hash, depth uint64) (index uint64, offset uint64) {
	hash = (hash + seeds[depth]) * seeds[depth]
	hash += (hash >> 32)
	index = hash & c.mask
	offset = ((hash&(counterDepth-1))*counterDepth + depth) * countBits
	return
}
