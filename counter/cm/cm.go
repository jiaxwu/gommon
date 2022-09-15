package cm

import (
	"math"

	"github.com/jiaxwu/gommon/hash"
	mmath "github.com/jiaxwu/gommon/math"
	"github.com/jiaxwu/gommon/mem"
	"golang.org/x/exp/constraints"
)

// Count-Min Sketch 计数器，原理类似于布隆过滤器，根据哈希映射到多个位置，然后在对应位置进行计数
// 读取时拿对应位置最小的
// 适合需要一个比较小的计数，而且不需要这个计数一定准确的情况
// 可以减少空间消耗
// https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.591.8351&rep=rep1&type=pdf
type Counter[T constraints.Unsigned] struct {
	counters    [][]T
	countersLen uint64       // 计数器长度
	hashs       []*hash.Hash // 哈希函数列表
	maxCount    T            // 最大计数值
}

// 创建一个计数器
// size：数据流大小
// errorRange：计数值误差范围（会超过真实计数值）
// errorRate：错误率
func New[T constraints.Unsigned](size uint64, errorRange T, errorRate float64) *Counter[T] {
	// 计数器长度
	countersLen := uint64(math.Ceil(math.E * float64(size) / float64(errorRange)))
	// 哈希个数
	hashsCnt := int(math.Ceil(math.Log(1 / errorRate)))
	hashs := make([]*hash.Hash, hashsCnt)
	counters := make([][]T, hashsCnt)
	for i := 0; i < hashsCnt; i++ {
		hashs[i] = hash.New()
		counters[i] = make([]T, countersLen)
	}
	return &Counter[T]{
		counters:    counters,
		countersLen: countersLen,
		hashs:       hashs,
		maxCount:    T(0) - 1,
	}
}

// 创建一个计数器
// size：数据流大小
// elements：不同元素数量
// errorRate：错误率
func NewWithElements[T constraints.Unsigned](size, elements uint64, errorRate float64) *Counter[T] {
	if elements > size {
		panic("too much elements")
	}
	errorRange := T(0) - 1
	if size/elements < uint64(errorRange) {
		errorRange = T(size / elements)
	}
	return New(size, errorRange, errorRate)
}

// 增加元素的计数
func (c *Counter[T]) Add(b []byte, val T) {
	for i, h := range c.hashs {
		index := h.Sum64(b) % c.countersLen
		if c.counters[i][index]+val <= c.counters[i][index] {
			c.counters[i][index] = c.maxCount
		} else {
			c.counters[i][index] += val
		}
	}
}

// 增加元素的计数
// 等同于Add(b, 1)
func (c *Counter[T]) Inc(b []byte) {
	c.Add(b, 1)
}

// 增加元素的计数
// 字符串类型
func (c *Counter[T]) AddString(s string, val T) {
	c.Add([]byte(s), val)
}

// 增加元素的计数
// 等同于Add(b, 1)
// 字符串类型
func (c *Counter[T]) IncString(s string) {
	c.Add([]byte(s), 1)
}

// 估算元素的计数
func (c *Counter[T]) Estimate(b []byte) T {
	minCount := c.maxCount
	for i, h := range c.hashs {
		index := h.Sum64(b) % c.countersLen
		count := c.counters[i][index]
		if count == 0 {
			return 0
		}
		minCount = mmath.Min(minCount, count)
	}
	return minCount
}

// 估算元素的计数
// 字符串类型
func (c *Counter[T]) EstimateString(s string) T {
	return c.Estimate([]byte(s))
}

// 计数衰减
// 如果factor为0则直接清空
func (c *Counter[T]) Attenuation(factor T) {
	for _, counter := range c.counters {
		if factor == 0 {
			mem.Memset(counter, 0)
		} else {
			for j := uint64(0); j < c.countersLen; j++ {
				counter[j] /= factor
			}
		}
	}
}

// 计数器数量
func (c *Counter[T]) Counters() uint64 {
	return c.countersLen
}

// 哈希函数数量
func (c *Counter[T]) Hashs() uint64 {
	return uint64(len(c.hashs))
}
