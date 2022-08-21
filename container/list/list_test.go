package list

import "testing"

func TestList(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)
	i := 1
	for !l.Empty() {
		first := l.Remove(l.Front())
		if i != first {
			t.Errorf("AddLast() = %v, want %v", first, i)
		}
		i++
	}
}
