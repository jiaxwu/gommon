package tinylfu

// 多少位表示一个计数
const counterDepth = 4

// 种子
var seeds = []uint64{0xc3a5c85c97cb3127, 0xb492b66fbe98f273, 0x9ae16a3b2f90404f, 0xcbf29ce484222325}

type CountMin struct {
	counters []uint64
	mask     uint64
	samples  int
	// Add的次数
	times int
}

// 计数，类似于布隆过滤器，根据哈希映射到多个位置，然后在对应位置进行计数
// 读取时拿对应位置最小的
func NewCountMin() {}

// 根据增加对应位置的计数
func (c *CountMin) Inc(hash uint64) {

}

// 返回counters下标
func (c *CountMin) index(hash, depth uint64) uint64 {
	hash = (hash + seeds[depth]) * seeds[depth]
	hash += (hash >> 32)
	return hash & c.mask
}

// 估算对应hash的计数
func (c *CountMin) Estimate(hash uint64) int {
	return 0
}
