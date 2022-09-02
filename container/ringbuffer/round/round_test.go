package round

import (
	"testing"

	"github.com/jiaxwu/gommon/container/ringbuffer/fix"
)

func TestRingAll(t *testing.T) {
	n := 10
	r := New[int](uint64(n))
	for i := 1; i <= n; i++ {
		r.Push(i)
	}

	v := r.Peek()
	if v != 1 {
		t.Errorf("Peek() = %v, want %v", v, 1)
	}

	i := 1
	for !r.Empty() {
		v := r.Pop()
		if i != v {
			t.Errorf("Pop() = %v, want %v", v, i)
		}
		i++
	}
}

func TestRingMAll(t *testing.T) {
	n := 10
	r := New[int](uint64(n))
	r.MPush(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	v := r.Peek()
	if v != 1 {
		t.Errorf("Peek() = %v, want %v", v, 1)
	}

	i := 1
	for !r.Empty() {
		v := r.MPop(1)
		if i != v[0] {
			t.Errorf("Pop() = %v, want %v", v, i)
		}
		i++
	}

	r.MPush(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	i = 1
	for !r.Empty() {
		v := r.MPop(1)
		if i != v[0] {
			t.Errorf("Pop() = %v, want %v", v, i)
		}
		i++
	}

	r.Reset()
	r.MPush(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	i = 1
	dst := make([]int, 1)
	for !r.Empty() {
		r.MPopCopy(dst)
		if i != dst[0] {
			t.Errorf("Pop() = %v, want %v", v, i)
		}
		i++
	}
}

const (
	RoundFixSize = 100000
)

func BenchmarkRoundPushPop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := New[int](RoundFixSize)
		for j := 0; j < RoundFixSize; j++ {
			r.Push(j)
		}
		for j := 0; j < RoundFixSize; j++ {
			r.Pop()
		}
	}
}

func BenchmarkFixPushPop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := fix.New[int](RoundFixSize)
		for j := 0; j < RoundFixSize; j++ {
			r.Push(j)
		}
		for j := 0; j < RoundFixSize; j++ {
			r.Pop()
		}
	}
}
