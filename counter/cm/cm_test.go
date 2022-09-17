package cm

import (
	"math"
	"strconv"
	"testing"
)

func TestCount(t *testing.T) {
	cm := New[uint8](1000, 10, 0.001)
	cm.AddString("10", 1)
	cm.AddString("51151", 1)
	cm.AddString("321", 1)
	cm.AddString("10", 1)
	cm.AddString("10", 1)
	cm.AddString("321", 1)
	if cm.EstimateString("10") != 3 {
		t.Errorf("want %v, but %d", 3, cm.EstimateString("10"))
	}
	if cm.EstimateString("321") != 2 {
		t.Errorf("want %v, but %d", 2, cm.EstimateString("321"))
	}
	if cm.EstimateString("51151") != 1 {
		t.Errorf("want %v, but %d", 1, cm.EstimateString("1"))
	}

	cm.AddString("10", 100)
	if cm.EstimateString("10") != 103 {
		t.Errorf("want %v, but %d", 103, cm.EstimateString("10"))
	}
	cm.AddString("10", 254)
	if cm.EstimateString("10") != 255 {
		t.Errorf("want %v, but %d", 255, cm.EstimateString("10"))
	}
	cm.AddString("5", 100)
	if cm.EstimateString("5") != 100 {
		t.Errorf("want %v, but %d", 100, cm.EstimateString("5"))
	}
	cm.AddString("1", 100)
	if cm.EstimateString("1") != 100 {
		t.Errorf("want %v, but %d", 100, cm.EstimateString("1"))
	}

	cm.Attenuation(2)
	if cm.EstimateString("10") != 127 {
		t.Errorf("want %v, but %d", 127, cm.EstimateString("10"))
	}
	if cm.EstimateString("5") != 50 {
		t.Errorf("want %v, but %d", 50, cm.EstimateString("5"))
	}
	if cm.EstimateString("1") != 50 {
		t.Errorf("want %v, but %d", 50, cm.EstimateString("1"))
	}
}

func TestExpectedErrorAndErrorRate(t *testing.T) {
	capacity := uint64(1000000)
	errorRange := uint32(1)
	errorRate := 0.001
	cm := New(capacity, errorRange, errorRate)
	// 添加计数值
	for i := uint64(0); i < capacity; i++ {
		cm.Add(i, 1)
	}
	// 评估
	errorCount := 0
	errorSum := 0
	for i := uint64(0); i < capacity; i++ {
		val := cm.Estimate(i)
		if val > 1+errorRange {
			errorCount++
			errorSum += int(val) - 1
		}
	}
	estimateErrorRate := float64(errorCount) / float64(capacity)
	estimateError := float64(errorSum) / math.Max(1, float64(errorCount))
	if estimateErrorRate > errorRate {
		t.Errorf("errorRate not accuracy %v", estimateErrorRate)
	}
	if estimateError > float64(errorRange) {
		t.Errorf("errorRange not accuracy %v", estimateError)
	}
}

func BenchmarkAddAndEstimateBytes(b *testing.B) {
	buf := make([]byte, 8192)
	for length := 1; length <= cap(buf); length *= 2 {
		b.Run(strconv.Itoa(length), func(b *testing.B) {
			f := New[uint8](uint64(b.N), 10, 0.0001)
			buf = buf[:length]
			b.SetBytes(int64(length))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f.AddBytes(buf, 1)
				f.EstimateBytes(buf)
			}
		})
	}
}

func BenchmarkAddAndEstimate(b *testing.B) {
	for length := 1; length <= 8192; length *= 2 {
		b.Run(strconv.Itoa(length), func(b *testing.B) {
			f := New[uint8](uint64(b.N), 10, 0.0001)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f.Add(uint64(i), 1)
				f.Estimate(uint64(i))
			}
		})
	}
}
