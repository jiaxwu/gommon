package bloom

import (
	"testing"
)

func TestNew(t *testing.T) {
	f := New(10, 0.01)
	f.AddHash(10)
	f.AddHash(44)
	f.AddHash(66)
	if !f.ContainsHash(10) {
		t.Errorf("want %v, but %v", true, f.ContainsHash(10))
	}
	if !f.ContainsHash(44) {
		t.Errorf("want %v, but %v", true, f.ContainsHash(10))
	}
	if !f.ContainsHash(66) {
		t.Errorf("want %v, but %v", true, f.ContainsHash(10))
	}
	if f.ContainsHash(55) {
		t.Errorf("want %v, but %v", false, f.ContainsHash(10))
	}
}
