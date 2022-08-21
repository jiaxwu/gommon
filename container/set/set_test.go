package set

import (
	"testing"
)

func TestSet(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	if !s.Contains(2) {
		t.Errorf("expected contains: %d, but not contains", 2)
	}

	if s.Contains(4) {
		t.Errorf("expected not contains: %d, but contains", 4)
	}

	s.Remove(2)

	if s.Contains(2) {
		t.Errorf("expected not contains: %d, but contains", 2)
	}
}
