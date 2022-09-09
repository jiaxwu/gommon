package qps

import (
	"testing"
	"time"
)

const windowCnt = 100

func TestAdd(t *testing.T) {
	q := New(windowCnt)
	addTimes := 100000
	for i := 0; i < addTimes; i++ {
		q.Add()
	}
	w := q.Get()
	if w.TotalCnt != int64(addTimes) {
		t.Errorf("totalCnt: %d, expected: %d", w.TotalCnt, addTimes)
	}
}

func TestAdd_SleepOneSecond(t *testing.T) {
	q := New(windowCnt)
	addTimes := 100000
	for i := 0; i < addTimes; i++ {
		q.Add()
	}
	time.Sleep(time.Second)
	w := q.Get()
	if w.TotalCnt != 0 {
		t.Errorf("totalCnt: %d, expected: %d", w.TotalCnt, 0)
	}
}

func BenchmarkAdd(b *testing.B) {
	q := New(windowCnt)
	for n := 0; n < b.N; n++ {
		q.Add()
	}
}

func BenchmarkAddSince(b *testing.B) {
	q := New(windowCnt)
	for n := 0; n < b.N; n++ {
		q.AddSince(time.Now())
	}
}

func BenchmarkAddUseTime(b *testing.B) {
	q := New(windowCnt)
	for n := 0; n < b.N; n++ {
		now := time.Now()
		q.AddUseTime(time.Since(now))
	}
}

func BenchmarkAddGet(b *testing.B) {
	q := New(windowCnt)
	for n := 0; n < b.N; n++ {
		q.Add()
		q.Get()
	}
}
