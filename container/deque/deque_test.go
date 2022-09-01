package deque

import (
	"testing"
)

func TestDequeAll(t *testing.T) {
	d := New[int]()
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)
	d.PushBack(4)
	i := 1
	for !d.Empty() {
		first := d.PopFront()
		if i != first {
			t.Errorf("AddLast() = %v, want %v", first, i)
		}
		i++
	}
}
