package pool

import (
	"bytes"
	"testing"
)

const (
	blocks    = 64
	blockSize = 1024
)

var block = make([]byte, blockSize)

func BenchmarkByte(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var b []byte
		for i := 0; i < blocks; i++ {
			b = append(b, block...)
		}
	}
}

func BenchmarkMake(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b := make([]byte, 0, blocks*blockSize)
		for i := 0; i < blocks; i++ {
			b = append(b, block...)
		}
	}
}

func BenchmarkBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b := bytes.NewBuffer(make([]byte, 0, blocks*blockSize))
		for i := 0; i < blocks; i++ {
			b.Write(block)
		}
	}
}

func BenchmarkBytePool(b *testing.B) {
	pool := NewBytePool(0, blocks*blockSize)
	for n := 0; n < b.N; n++ {
		b := pool.Get()
		for i := 0; i < blocks; i++ {
			b = append(b, block...)
		}
		pool.Put(b)
	}
}

func BenchmarkBufferPool(b *testing.B) {
	pool := NewBufferPool(0, blocks*blockSize)
	for n := 0; n < b.N; n++ {
		b := pool.Get()
		for i := 0; i < blocks; i++ {
			b.Write(block)
		}
		pool.Put(b)
	}
}
