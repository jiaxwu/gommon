package deque

import (
	"testing"
)

func TestDeque_AddLast(t *testing.T) {
	d := New[int]()
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)
	d.PushBack(4)
	i := 1
	for !d.Empty() {
		first := d.RemoveFront()
		if i != first {
			t.Errorf("AddLast() = %v, want %v", first, i)
		}
		i++
	}
}
