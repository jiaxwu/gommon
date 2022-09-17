package cm

import (
	"hash/fnv"
	"math"
	"math/rand"
	"time"

	mmath "github.com/jiaxwu/gommon/math"
	"github.com/jiaxwu/gommon/mem"
)

const (
	// 计数器位数
	counter4Bits = 4
	// 最大计数值
	counter4MaxVal = 1<<counter4Bits - 1
)

// 4bit 版本 Count-Min Sketch 计数器
type Counter4 struct {
	counters   [][]uint64
	counterCnt uint64   // 计数器长度
	seeds      []uint64 // 哈希种子
}

// 创建一个计数器
// size：数据流大小
// errorRange：计数值误差范围（会超过真实计数值）
// errorRate：错误率
func New4(size uint64, errorRange uint8, errorRate float64) *Counter4 {
	if errorRange > counter4MaxVal {
		panic("too large errorRange")
	}
	// 计数器长度
	counterCnt := uint64(math.Ceil(math.E / (float64(errorRange) / float64(size)) / (64 / counter4Bits)))
	// 哈希个数
	seedCnt := int(math.Ceil(math.Log(1 / errorRate)))
	seeds := make([]uint64, seedCnt)
	counters := make([][]uint64, seedCnt)
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < seedCnt; i++ {
		seeds[i] = source.Uint64()
		counters[i] = make([]uint64, counterCnt)
	}
	return &Counter4{
		counters:   counters,
		counterCnt: counterCnt,
		seeds:      seeds,
	}
}

// 创建一个计数器
// size：数据流大小
// elements：不同元素数量
// errorRate：错误率
func New4WithElements(size, elements uint64, errorRate float64) *Counter4 {
	if elements > size {
		panic("too much elements")
	}
	errorRange := uint8(counter4MaxVal)
	if size/elements < uint64(errorRange) {
		errorRange = uint8(size / elements)
	}
	return New4(size, errorRange, errorRate)
}

// 增加元素的计数
func (c *Counter4) Add(h uint64, val uint8) {
	for i, seed := range c.seeds {
		index, offset := c.pos(h, seed)
		count := c.getCount(c.counters[i], index, offset)
		count += uint64(val)
		if count > counter4MaxVal {
			count = counter4MaxVal
		}
		c.setCount(c.counters[i], index, offset, count)
	}
}

// 增加元素的计数
func (c *Counter4) AddBytes(b []byte, val uint8) {
	c.Add(c.hash(b), val)
}

// 增加元素的计数
// 字符串类型
func (c *Counter4) AddString(s string, val uint8) {
	c.AddBytes([]byte(s), val)
}

// 估算元素的计数
func (c *Counter4) Estimate(h uint64) uint8 {
	minCount := uint8(counter4MaxVal)
	for i, seed := range c.seeds {
		index, offset := c.pos(h, seed)
		count := c.getCount(c.counters[i], index, offset)
		if count == 0 {
			return 0
		}
		minCount = mmath.Min(minCount, uint8(count))
	}
	return minCount
}

// 估算元素的计数
func (c *Counter4) EstimateBytes(b []byte) uint8 {
	return c.Estimate(c.hash(b))
}

// 估算元素的计数
// 字符串类型
func (c *Counter4) EstimateString(s string) uint8 {
	return c.EstimateBytes([]byte(s))
}

// 计数衰减
// 如果factor为0则直接清空
func (c *Counter4) Attenuation(factor uint8) {
	for _, counter := range c.counters {
		if factor == 0 || factor > counter4MaxVal {
			mem.Memset(counter, 0)
		} else {
			for index := uint64(0); index < c.counterCnt; index++ {
				for offset := uint64(0); offset < 64; offset += counter4Bits {
					count := c.getCount(counter, index, offset) / uint64(factor)
					c.setCount(counter, index, offset, count)
				}
			}
		}
	}
}

// 计数器数量
func (c *Counter4) Counters() uint64 {
	return c.counterCnt * (64 / counter4Bits)
}

// 哈希函数数量
func (c *Counter4) Hashs() uint64 {
	return uint64(len(c.seeds))
}

// 返回位置
// 也就是index和offset
func (c *Counter4) pos(h, seed uint64) (uint64, uint64) {
	// 哈希值
	hashValue := seed ^ h
	// 计数器下标
	index := hashValue % c.counterCnt
	// 计数器在64位里面的偏移
	offset := (hashValue & counter4MaxVal) * counter4Bits
	return index, offset
}

// 获取计数值
func (c *Counter4) getCount(counter []uint64, index, offset uint64) uint64 {
	return (counter[index] >> offset) & uint64(counter4MaxVal)
}

// 设置计数值
func (c *Counter4) setCount(counter []uint64, index, offset, count uint64) {
	counter[index] = (counter[index] &^ (counter4MaxVal << offset)) | (count << offset)
}

// 计算哈希值
func (c *Counter4) hash(b []byte) uint64 {
	f := fnv.New64()
	f.Write(b)
	return f.Sum64()
}
