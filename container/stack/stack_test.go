package deque

import (
	"testing"
)

func TestStack_Push(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	i := 3
	for !s.Empty() {
		last := s.Pop()
		if i != last {
			t.Errorf("Push() = %v, want %v", last, i)
		}
		i--
	}
}
