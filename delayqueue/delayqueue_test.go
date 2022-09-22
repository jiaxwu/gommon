package delayqueue

import (
	"context"
	"testing"
	"time"
)

func TestDelayQueue(t *testing.T) {
	times := []int{1, 2, 3}
	q := New[int, int]()
	for _, t := range times {
		q.Push(t, t, time.Microsecond*time.Duration(t))
	}

	for _, time := range times {
		entry, ok := q.Take(context.Background())
		if !ok {
			t.Errorf("want %v, but %v", true, ok)
		}
		if entry.Key != time || entry.Value != time {
			t.Errorf("want %v, but %v", time, entry.Key)
		}
	}
}

func TestRemove(t *testing.T) {
	times := []int{1, 2, 3}
	q := New[int, int]()
	for _, t := range times {
		q.Push(t, t, time.Microsecond*time.Duration(t))
	}

	q.Remove(2)

	entry, ok := q.Take(context.Background())
	if !ok {
		t.Errorf("want %v, but %v", true, ok)
	}
	if entry.Key == 2 || entry.Value == 2 {
		t.Errorf("invalid %v", 2)
	}

	entry, ok = q.Take(context.Background())
	if !ok {
		t.Errorf("want %v, but %v", true, ok)
	}
	if entry.Key == 2 || entry.Value == 2 {
		t.Errorf("invalid %v", 2)
	}
}

func Benchmark(b *testing.B) {
	q := New[int, int]()
	for i := 0; i < b.N; i++ {
		q.Push(i, i, 0)
	}
	ch := q.Channel(context.Background(), 10)
	for i := 0; i < b.N; i++ {
		_, ok := <-ch
		if !ok {
			b.Errorf("want %v, but %v", true, ok)
		}
		// if entry.Key != i || entry.Value != i {
		// 	b.Errorf("want %v, but %v", i, entry.Key)
		// }
	}
}
