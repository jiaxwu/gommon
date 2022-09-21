package heap

import (
	"testing"
)

func TestRemovableHeapAll(t *testing.T) {
	h := NewRemovableHeap(func(e1 Entry[int, int], e2 Entry[int, int]) bool {
		return e1.Value > e2.Value
	})
	h.Push(1, 5)
	h.Push(2, 6)
	h.Push(3, 3)
	h.Push(4, 7)
	h.Push(5, 2)
	h.Push(6, 4)
	h.Push(7, 8)
	h.Push(8, 9)
	h.Push(9, 1)

	e := h.Peek()
	if e.Value != 9 {
		t.Errorf("Peek() = %v, want %v", e.Value, 9)
	}

	i := 9
	for !h.Empty() {
		e := h.Pop()
		if i != e.Value {
			t.Errorf("Pop() = %v, want %v", e.Value, i)
		}
		i--
	}
}

func TestRemovableHeapRemove(t *testing.T) {
	h := NewRemovableHeap(func(e1 Entry[int, int], e2 Entry[int, int]) bool {
		return e1.Value > e2.Value
	})
	h.Push(1, 5)
	h.Push(2, 6)
	h.Push(3, 3)
	h.Push(4, 7)
	h.Push(5, 2)
	h.Push(6, 4)
	h.Push(7, 8)
	h.Push(8, 9)
	h.Push(9, 1)

	e := h.Peek()
	if e.Value != 9 {
		t.Errorf("Peek() = %v, want %v", e.Value, 9)
	}

	v, ok := h.Get(5)
	if !ok || v != 2 {
		t.Errorf("Get() = %v, want %v", v, 2)
	}

	h.Remove(5)

	v, ok = h.Get(5)
	if ok || v != 0 {
		t.Errorf("Get() = %v, want %v", v, 0)
	}

	i := 9
	for !h.Empty() {
		e := h.Pop()
		if i != e.Value {
			t.Errorf("Pop() = %v, want %v", e.Value, i)
		}
		if i == 3 {
			i--
		}
		i--
	}
}

func FuzzRemovableHeap(f *testing.F) {
	seeds := [][]int{{1, 2}, {4, 6}, {3, 2}}
	for _, seed := range seeds {
		f.Add(seed[0], seed[1])
	}
	h := NewRemovableHeap(func(e1 Entry[int, int], e2 Entry[int, int]) bool {
		return e1.Value > e2.Value
	})
	m := map[int]bool{}
	f.Fuzz(func(t *testing.T, key, value int) {
		if m[key] {
			h.Remove(key)
			delete(m, key)
			return
		}
		m[key] = true
		h.Push(key, value)
	})
}
