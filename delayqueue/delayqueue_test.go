package delayqueue

import (
	"context"
	"testing"
	"time"
)

func TestDelayQueue(t *testing.T) {
	times := []int{1, 2, 3}
	q := New[int]()
	for _, t := range times {
		q.Push(t, time.Microsecond*time.Duration(t))
	}

	for _, time := range times {
		value, ok := q.Take(context.Background())
		if !ok {
			t.Errorf("want %v, but %v", true, ok)
		}
		if value != time {
			t.Errorf("want %v, but %v", time, value)
		}
	}
}

func BenchmarkPushAndTake(b *testing.B) {
	q := New[int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Push(i, time.Duration(i))
	}
	b.StopTimer()
	time.Sleep(time.Duration(b.N))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, ok := q.Take(context.Background())
		if !ok {
			b.Errorf("want %v, but %v", true, ok)
		}
	}
}
