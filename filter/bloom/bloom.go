package bloom

import (
	"math"

	"github.com/jiaxwu/gommon/hash"
)

// uint64的位数
const uint64Bits = 64

// 布隆过滤器
// https://llimllib.github.io/bloomfilter-tutorial/
// https://github.com/bits-and-blooms/bloom/blob/master/bloom.go
type Filter struct {
	bits    []uint64     // bit数组
	bitsCnt uint64       // bit位数
	hashs   []*hash.Hash // 不同哈希函数
}

// capacity：容量
// falsePositiveRate：误判率
func New(capacity uint64, falsePositiveRate float64) *Filter {
	// bit数量
	factor := -math.Log(falsePositiveRate) / (math.Ln2 * math.Ln2)
	bitsCnt := uint64(math.Ceil(float64(capacity) * factor))
	// 这里扩大到最后一个uint64大小，避免浪费
	bitsCnt = (bitsCnt + uint64Bits - 1) / uint64Bits * uint64Bits

	// 哈希函数数量
	hashsCnt := int(math.Ceil(math.Ln2 * float64(bitsCnt) / float64(capacity)))
	hashs := make([]*hash.Hash, hashsCnt)
	for i := 0; i < hashsCnt; i++ {
		hashs[i] = hash.New()
	}

	return &Filter{
		bits:    make([]uint64, bitsCnt/uint64Bits),
		bitsCnt: bitsCnt,
		hashs:   hashs,
	}
}

// 添加元素
func (f *Filter) Add(b []byte) {
	for _, h := range f.hashs {
		index, offset := f.pos(h, b)
		f.bits[index] |= 1 << offset
	}
}

// 添加元素
// 字符串类型
func (f *Filter) AddString(s string) {
	f.Add([]byte(s))
}

// 元素是否存在
// true表示可能存在
func (f *Filter) Contains(b []byte) bool {
	for _, h := range f.hashs {
		index, offset := f.pos(h, b)
		mask := uint64(1) << offset
		// 判断这一位是否位1
		if (f.bits[index] & mask) != mask {
			return false
		}
	}
	return true
}

// 元素是否存在
// 字符串类型
func (f *Filter) ContainsString(s string) bool {
	return f.Contains([]byte(s))
}

// 清空过滤器
func (f *Filter) Clear() {
	for i := range f.bits {
		f.bits[i] = 0
	}
}

// 获取对应元素下标和偏移
func (f *Filter) pos(h *hash.Hash, b []byte) (uint64, uint64) {
	hashValue := h.Sum64(b)
	// 按照位计算的偏移
	bitsIndex := hashValue % f.bitsCnt
	// 因为一个元素64位，因此需要转换
	index := bitsIndex / uint64Bits
	// 在一个元素里面的偏移
	offset := bitsIndex % uint64Bits
	return index, offset
}
