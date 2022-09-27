package heap

import (
	"testing"
)

func TestHeapNAll(t *testing.T) {
	h := NewN(4, func(e1 int, e2 int) bool {
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

	v := h.Peek()
	if v != 9 {
		t.Errorf("Peek() = %v, want %v", v, 9)
	}

	i := 9
	for !h.Empty() {
		v := h.Pop()
		if i != v {
			t.Errorf("Pop() = %v, want %v", v, i)
		}
		i--
	}
}
