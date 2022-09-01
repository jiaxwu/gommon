package priority_queue

import (
	"testing"
)

func TestPriorityQueueAll(t *testing.T) {
	q := New(nil, func(e1 int, e2 int) bool {
		return e1 > e2
	})
	q.Push(5)
	q.Push(6)
	q.Push(3)
	q.Push(7)
	q.Push(2)
	q.Push(4)
	q.Push(8)
	q.Push(9)
	q.Push(1)

	v := q.Peek()
	if v != 9 {
		t.Errorf("Peek() = %v, want %v", v, 9)
	}

	i := 9
	for !q.Empty() {
		v := q.Pop()
		if i != v {
			t.Errorf("Pop() = %v, want %v", v, i)
		}
		i--
	}
}
