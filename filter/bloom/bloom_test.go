package bloom

import (
	"encoding/binary"
	"hash/fnv"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	f := New(10, 0.01)
	f.AddString("10")
	f.AddString("44")
	f.AddString("66")
	if !f.ContainsString("10") {
		t.Errorf("want %v, but %v", true, f.ContainsString("10"))
	}
	if !f.ContainsString("44") {
		t.Errorf("want %v, but %v", true, f.ContainsString("10"))
	}
	if !f.ContainsString("66") {
		t.Errorf("want %v, but %v", true, f.ContainsString("10"))
	}
	if f.ContainsString("55") {
		t.Errorf("want %v, but %v", false, f.ContainsString("10"))
	}
}

func BenchmarkAddAndContains(b *testing.B) {
	buf := make([]byte, 8192)
	for length := 1; length <= cap(buf); length *= 2 {
		b.Run(strconv.Itoa(length), func(b *testing.B) {
			f := New(uint64(b.N), 0.0001)
			buf = buf[:length]
			b.SetBytes(int64(length))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f.AddBytes(buf)
				f.ContainsBytes(buf)
			}
		})
	}
}

func TestFalsePositiveRate(t *testing.T) {
	capacity := uint64(10000000)
	rounds := uint64(10000000)
	falsePositiveRate := 0.01
	f := New(capacity, falsePositiveRate)
	// 加入过滤器一些元素
	for i := uint64(0); i < capacity; i++ {
		h := fnv.New64()
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, i)
		h.Write(b)
		f.Add(h.Sum64())
	}
	// 查询不存在的元素，计算错误率
	falsePositiveCount := 0
	for i := uint64(0); i < rounds; i++ {
		// 加上容量保证这个元素一定不是之前加入过滤器的
		h := fnv.New64()
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, i+capacity+1)
		h.Write(b)
		if f.Contains(h.Sum64()) {
			falsePositiveCount++
		}
	}
	t.Log(falsePositiveCount)
	fpRate := float64(falsePositiveCount) / (float64(rounds))
	if !(fpRate >= falsePositiveRate*(0.9) && fpRate <= falsePositiveRate*(1.1)) {
		t.Errorf("fpRate not accuracy %v", fpRate)
	}
}
