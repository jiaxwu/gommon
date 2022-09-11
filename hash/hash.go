package hash

import (
	"hash/maphash"
)

// 哈希函数
// 非线程安全，业务请加锁
type Hasher struct {
	h *maphash.Hash
}

func New() *Hasher {
	h := &Hasher{
		h: &maphash.Hash{},
	}
	h.h.SetSeed(maphash.MakeSeed())
	return h
}

// 计算哈希值
func (h *Hasher) Sum64(b []byte) uint64 {
	h.h.Reset()
	h.h.Write(b)
	return h.h.Sum64()
}

// 计算哈希值
func (h *Hasher) Sum64String(s string) uint64 {
	return h.Sum64([]byte(s))
}
