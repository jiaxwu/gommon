package level

import (
	"math/rand"
	"sync/atomic"
	"testing"

	"github.com/jiaxwu/gommon/math"

	"github.com/jiaxwu/gommon/pool"
)

const (
	blocks    = 64
	blockSize = 1024
)

var block = make([]byte, blockSize)

func BenchmarkBytePoolRandomBlcoks(b *testing.B) {
	var eqCap int64
	pool := pool.NewBytePool(0, 0)
	for n := 0; n < b.N; n++ {
		blocks := rand.Intn(blocks) + 1
		needCap := blocks * blockSize
		b := pool.Get()
		if cap(b) >= needCap && cap(b) <= needCap*2 {
			atomic.AddInt64(&eqCap, 1)
		}
		for i := 0; i < blocks; i++ {
			b = append(b, block...)
		}
		pool.Put(b)
	}
	b.Logf("eq cap cnt: %d", eqCap)
}

func BenchmarkLevelBytePoolRandomBlcoks(b *testing.B) {
	var eqCap int64
	var pool = NewBytePool(blockSize, math.Log2(blocks)+1)
	for n := 0; n < b.N; n++ {
		blocks := rand.Intn(blocks) + 1
		needCap := blocks * blockSize
		b := pool.Get(needCap)
		if cap(b) >= needCap && cap(b) <= needCap*2 {
			atomic.AddInt64(&eqCap, 1)
		}
		for i := 0; i < blocks; i++ {
			b = append(b, block...)
		}
		pool.Put(b)
	}
	b.Logf("eq cap cnt: %d", eqCap)
}

func BenchmarkLevelBufferPoolRandomBlcoks(b *testing.B) {
	var eqCap int64
	var pool = NewBufferPool(blockSize, math.Log2(blocks)+1)
	for n := 0; n < b.N; n++ {
		blocks := rand.Intn(blocks) + 1
		needCap := blocks * blockSize
		b := pool.Get(needCap)
		if b.Cap() >= needCap && b.Cap() <= needCap*2 {
			atomic.AddInt64(&eqCap, 1)
		}
		for i := 0; i < blocks; i++ {
			b.Write(block)
		}
		pool.Put(b)
	}
	b.Logf("eq cap cnt: %d", eqCap)
}
