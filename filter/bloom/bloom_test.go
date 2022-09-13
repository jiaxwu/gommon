package bloom

import (
	"encoding/binary"
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
				f.Add(buf)
				f.Contains(buf)
			}
		})
	}
}

func TestFalsePositiveRate(t *testing.T) {
	capacity := uint64(100000)
	rounds := uint32(100000)
	falsePositiveRate := 0.01
	f := New(capacity, falsePositiveRate)
	// 加入过滤器一些元素
	item := make([]byte, 4)
	for i := uint32(0); i < uint32(capacity); i++ {
		binary.BigEndian.PutUint32(item, i)
		f.Add(item)
	}
	// 查询不存在的元素，计算错误率
	falsePositiveCount := 0
	for i := uint32(0); i < rounds; i++ {
		// 加上容量保证这个元素一定不是之前加入过滤器的
		binary.BigEndian.PutUint32(item, i+uint32(capacity)+1)
		if f.Contains(item) {
			falsePositiveCount++
		}
	}
	fpRate := float64(falsePositiveCount) / (float64(rounds))
	if !(fpRate >= falsePositiveRate-0.001 && fpRate <= falsePositiveRate+0.001) {
		t.Errorf("fpRate not accuracy %v", fpRate)
	}
}
