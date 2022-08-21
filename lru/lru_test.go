package lru

import (
	"testing"
)

func TestLRU_Put(t *testing.T) {
	l := New[string, int](3)
	l.Put("11", 5)
	l.Put("22", 6)
	l.Put("33", 7)
	l.Get("11")
	l.Put("44", 8)

	get22, ok := l.Get("22")
	if get22 != 0 || ok {
		t.Errorf("Put() = %v, want %v", ok, false)
	}
}
