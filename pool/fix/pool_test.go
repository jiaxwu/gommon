package fix

import "testing"

const (
	blocks    = 64
	blockSize = 1024
)

var block = make([]byte, blockSize)

func BenchmarkBytePool(b *testing.B) {
	pool := NewBytePool(16, 0, blocks*blockSize)
	for n := 0; n < b.N; n++ {
		b := pool.Get()
		for i := 0; i < blocks; i++ {
			b = append(b, block...)
		}
		pool.Put(b)
	}
}

func BenchmarkBufferPool(b *testing.B) {
	pool := NewBufferPool(16, 0, blocks*blockSize)
	for n := 0; n < b.N; n++ {
		b := pool.Get()
		for i := 0; i < blocks; i++ {
			b.Write(block)
		}
		pool.Put(b)
	}
}
