package meap

import (
	"testing"

	"github.com/jiaxwu/gommon/container/heap"
	"github.com/jiaxwu/gommon/container/meap"
)

func TestAll(t *testing.T) {
	h := New(func(e1, e2 int) bool {
		return e1 > e2
	})
	h.Push(5)
	h.Push(6)
	h.Push(3)
	h.Push(7)
	h.Push(2)
	h.Push(4)
	h.Push(8)
	h.Push(9)
	h.Push(1)

	e := h.Peek()
	if e != 9 {
		t.Errorf("Peek() = %v, want %v", e, 9)
	}

	i := 9
	for !h.Empty() {
		e := h.Pop()
		if i != e {
			t.Errorf("Pop() = %v, want %v", e, i)
		}
		i--
	}
}

func TestRemove(t *testing.T) {
	h := New(func(e1, e2 int) bool {
		return e1 > e2
	})
	entry := h.Push(5)
	h.Push(6)
	h.Push(3)
	h.Push(7)
	h.Push(2)
	h.Push(4)
	h.Push(8)
	h.Push(9)
	h.Push(1)

	e := h.Peek()
	if e != 9 {
		t.Errorf("Peek() = %v, want %v", e, 9)
	}

	h.Remove(entry)

	i := 9
	for !h.Empty() {
		e := h.Pop()
		if i != e {
			t.Errorf("Pop() = %v, want %v", e, i)
		}
		if i == 6 {
			i--
		}
		i--
	}
}

func Fuzz(f *testing.F) {
	seeds := []int{1, 2, 4, 6, 3, 2}
	for _, seed := range seeds {
		f.Add(seed)
	}
	h := New(func(e1, e2 int) bool {
		return e1 > e2
	})
	m := map[int]*Entry[int]{}
	f.Fuzz(func(t *testing.T, value int) {
		if m[value] != nil {
			h.Remove(m[value])
			delete(m, value)
			return
		}
		m[value] = h.Push(value)
	})
}

func BenchmarkHeapPushAndPop(b *testing.B) {
	q := heap.New(nil, func(e1, e2 int) bool {
		return e1 < e2
	})
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}

	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}

func BenchmarkReapPushAndPop(b *testing.B) {
	q := New(func(e1, e2 int) bool {
		return e1 < e2
	})
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}

	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}

func BenchmarkMeapPushAndPop(b *testing.B) {
	q := meap.New(func(e1, e2 meap.Entry[int, int]) bool {
		return e1.Value < e2.Value
	})
	for i := 0; i < b.N; i++ {
		q.Push(i, i)
	}

	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
