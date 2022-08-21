package priority_queue

import (
	"testing"
)

func TestPriorityQueue_Add(t *testing.T) {
	h := New(nil, func(e1 int, e2 int) bool {
		return e1 > e2
	})
	h.Add(5)
	h.Add(6)
	h.Add(3)
	h.Add(7)
	h.Add(2)
	h.Add(4)
	h.Add(8)
	h.Add(9)
	h.Add(1)

	i := 9
	for !h.Empty() {
		v := h.Remove()
		if i != v {
			t.Errorf("Add() = %v, want %v", v, i)
		}
		i--
	}
}
