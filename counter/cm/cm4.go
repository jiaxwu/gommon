package cm

import (
	"math"

	"github.com/jiaxwu/gommon/hash"
	mmath "github.com/jiaxwu/gommon/math"
	"github.com/jiaxwu/gommon/mem"
)

const (
	// 计数器位数
	counter4Bits = 4
	// 最大计数值
	counter4MaxCount = 1<<counter4Bits - 1
)

// 4bit 版本 Count-Min Sketch 计数器
type Counter4 struct {
	counters    [][]uint64
	countersLen uint64       // 计数器长度
	hashs       []*hash.Hash // 哈希函数列表
}

// 创建一个计数器
// size：数据流大小
// errorRange：计数值误差范围（会超过真实计数值）
// errorRate：错误率
func New4(size uint64, errorRange uint8, errorRate float64) *Counter4 {
	if errorRange > counter4MaxCount {
		panic("too large errorRange")
	}
	// 计数器长度
	countersLen := uint64(math.Ceil(math.E / (float64(errorRange) / float64(size)) / (64 / counter4Bits)))
	// 哈希个数
	hashsCnt := int(math.Ceil(math.Log(1 / errorRate)))
	hashs := make([]*hash.Hash, hashsCnt)
	counters := make([][]uint64, hashsCnt)
	for i := 0; i < hashsCnt; i++ {
		hashs[i] = hash.New()
		counters[i] = make([]uint64, countersLen)
	}
	return &Counter4{
		counters:    counters,
		countersLen: countersLen,
		hashs:       hashs,
	}
}

// 增加元素的计数
func (c *Counter4) Add(b []byte, val uint8) {
	for i, h := range c.hashs {
		index, offset := c.pos(h, b)
		count := c.getCount(c.counters[i], index, offset)
		count += uint64(val)
		if count > counter4MaxCount {
			count = counter4MaxCount
		}
		c.setCount(c.counters[i], index, offset, count)
	}
}

// 增加元素的计数
// 等同于Add(b, 1)
func (c *Counter4) Inc(b []byte) {
	c.Add(b, 1)
}

// 增加元素的计数
// 字符串类型
func (c *Counter4) AddString(s string, val uint8) {
	c.Add([]byte(s), val)
}

// 增加元素的计数
// 等同于Add(b, 1)
// 字符串类型
func (c *Counter4) IncString(s string) {
	c.Add([]byte(s), 1)
}

// 估算元素的计数
func (c *Counter4) Estimate(b []byte) uint8 {
	minCount := uint8(counter4MaxCount)
	for i, h := range c.hashs {
		index, offset := c.pos(h, b)
		count := c.getCount(c.counters[i], index, offset)
		if count == 0 {
			return 0
		}
		minCount = mmath.Min(minCount, uint8(count))
	}
	return minCount
}

// 估算元素的计数
// 字符串类型
func (c *Counter4) EstimateString(s string) uint8 {
	return c.Estimate([]byte(s))
}

// 计数衰减
// 如果factor为0则直接清空
func (c *Counter4) Attenuation(factor uint8) {
	for _, counter := range c.counters {
		if factor == 0 || factor > counter4MaxCount {
			mem.Memset(counter, 0)
		} else {
			for index := uint64(0); index < c.countersLen; index++ {
				for offset := uint64(0); offset < 64; offset += counter4Bits {
					count := c.getCount(counter, index, offset) / uint64(factor)
					c.setCount(counter, index, offset, count)
				}
			}
		}
	}
}

// 计数器长度
func (c *Counter4) Len() uint64 {
	return c.countersLen * (64 / counter4Bits)
}

// 返回位置
// 也就是index和offset
func (c *Counter4) pos(h *hash.Hash, b []byte) (uint64, uint64) {
	// 哈希值
	hashValue := h.Sum64(b)
	// 计数器下标
	index := hashValue % c.countersLen
	// 计数器在64位里面的偏移
	offset := (hashValue & counter4MaxCount) * counter4Bits
	return index, offset
}

// 获取计数值
func (c *Counter4) getCount(counter []uint64, index, offset uint64) uint64 {
	return (counter[index] >> offset) & uint64(counter4MaxCount)
}

// 设置计数值
func (c *Counter4) setCount(counter []uint64, index, offset, count uint64) {
	counter[index] = (counter[index] &^ (counter4MaxCount << offset)) | (count << offset)
}
