package meap

import (
	"testing"
)

func TestAll(t *testing.T) {
	h := New(func(e1, e2 Entry[int, int]) bool {
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

	h.Push(9, 1)
	h.Push(9, 2)
	e = h.Peek()
	if e.Value != 2 {
		t.Errorf("Peek() = %v, want %v", e.Value, 2)
	}
}

func TestRemove(t *testing.T) {
	h := New(func(e1, e2 Entry[int, int]) bool {
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

func Fuzz(f *testing.F) {
	seeds := [][]int{{1, 2}, {4, 6}, {3, 2}}
	for _, seed := range seeds {
		f.Add(seed[0], seed[1])
	}
	h := New(func(e1, e2 Entry[int, int]) bool {
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
