package bloom

import (
	"math"

	"github.com/jiaxwu/gommon/hash"

	mmath "github.com/jiaxwu/gommon/math"
)

// 布隆过滤器
// https://llimllib.github.io/bloomfilter-tutorial/
// https://github.com/bits-and-blooms/bloom/blob/master/bloom.go
type Filter struct {
	bits     []uint64     // bit数组
	bitsMask uint64       // bit数组掩码，也就是bits数组长度-1，用于快速取模
	hashs    []*hash.Hash // 不同哈希函数
}

// capacity：容量
// falsePositiveRate：误判率
func New(capacity uint64, falsePositiveRate float64) *Filter {
	// bit数量
	ln2 := math.Log(2.0)
	factor := -math.Log(falsePositiveRate) / (ln2 * ln2)
	bits := mmath.RoundUpPowOf2(uint64(float64(capacity) * factor))
	if bits == 0 {
		bits = 1
	}
	bitsMask := bits - 1

	// 哈希函数数量
	hashsLen := int(ln2 * float64(bits) / float64(capacity))
	if hashsLen < 1 {
		hashsLen = 1
	}
	hashs := make([]*hash.Hash, hashsLen)
	for i := 0; i < hashsLen; i++ {
		hashs[i] = hash.New()
	}

	return &Filter{
		bits:     make([]uint64, (bits+63)/64),
		bitsMask: bitsMask,
		hashs:    hashs,
	}
}

// 添加
func (f *Filter) Add(b []byte) {
	for _, h := range f.hashs {
		hashValue := h.Sum64(b)
		f.set(hashValue & f.bitsMask)
	}
}

func (f *Filter) AddString(s string) {
	f.Add([]byte(s))
}

// 如果可能存在则返回true
func (f *Filter) Contains(b []byte) bool {
	exists := true
	for _, h := range f.hashs {
		hashValue := h.Sum64(b)
		exists = f.get(hashValue&f.bitsMask) && exists
	}
	return exists
}

func (f *Filter) ContainsString(s string) bool {
	return f.Contains([]byte(s))
}

// 设置对应下标的值
// 如果对应下标已经为1则返回true
func (f *Filter) set(index uint64) {
	idx := index / 64
	shift := index % 64
	f.bits[idx] |= 1 << shift
}

// 获取对应下标的值
// 如果1返回true
func (f *Filter) get(index uint64) bool {
	idx := index / 64
	shift := index % 64
	val := f.bits[idx]
	mask := uint64(1) << shift
	return (val&mask)>>shift == 1
}

// 清空过滤器
func (f *Filter) Reset() {
	for i := range f.bits {
		f.bits[i] = 0
	}
}
