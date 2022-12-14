package hash

import (
	"hash/maphash"
)

// 哈希函数
// 非线程安全，业务请加锁
// 也就是对maphash的包装
type Hash struct {
	h *maphash.Hash
}

func New() *Hash {
	h := &Hash{
		h: &maphash.Hash{},
	}
	h.h.SetSeed(maphash.MakeSeed())
	return h
}

// 计算哈希值
func (h *Hash) Sum64(b []byte) uint64 {
	h.h.Reset()
	h.h.Write(b)
	return h.h.Sum64()
}

// 计算哈希值
func (h *Hash) Sum64String(s string) uint64 {
	return h.Sum64([]byte(s))
}
