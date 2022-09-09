package hash

import (
	"hash/maphash"
	"unsafe"
)

// 泛型哈希函数
// 非线程安全，业务请加锁
type Hasher[T any] struct {
	h    *maphash.Hash
	size int
}

func New[T any]() *Hasher[T] {
	var t T
	h := &Hasher[T]{
		h:    &maphash.Hash{},
		size: int(unsafe.Sizeof(t)),
	}
	h.h.SetSeed(maphash.MakeSeed())
	return h
}

// 计算哈希值
func (h *Hasher[T]) Hash(t T) uint64 {
	b := *(*[]byte)(unsafe.Pointer(&struct {
		data unsafe.Pointer
		len  int
	}{unsafe.Pointer(&t), h.size}))
	h.h.Reset()
	h.h.Write(b)
	return h.h.Sum64()
}
