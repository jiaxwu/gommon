package hash

import (
	"testing"
)

type Key struct {
	S string
	B int
}

func TestStruct(t *testing.T) {
	h := New[Key]()
	a := h.Hash(Key{
		S: "ab",
		B: 6,
	})
	b := h.Hash(Key{
		S: "ab",
		B: 6,
	})
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func TestPointer(t *testing.T) {
	h := New[*Key]()
	k := &Key{
		S: "ab",
		B: 6,
	}
	a := h.Hash(k)
	b := h.Hash(k)
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func TestInt(t *testing.T) {
	h := New[int]()
	a := h.Hash(32)
	b := h.Hash(32)
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func TestInt8(t *testing.T) {
	h := New[int8]()
	a := h.Hash(32)
	b := h.Hash(32)
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func TestInt16(t *testing.T) {
	h := New[int16]()
	a := h.Hash(32)
	b := h.Hash(32)
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func TestInt32(t *testing.T) {
	h := New[int32]()
	a := h.Hash(32)
	b := h.Hash(32)
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func TestInt64(t *testing.T) {
	h := New[int64]()
	a := h.Hash(32)
	b := h.Hash(32)
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func TestFloat32(t *testing.T) {
	h := New[float32]()
	a := h.Hash(32.1)
	b := h.Hash(32.1)
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func TestFloat64(t *testing.T) {
	h := New[float64]()
	a := h.Hash(32.2)
	b := h.Hash(32.2)
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func TestString(t *testing.T) {
	h := New[string]()
	a := h.Hash("bz")
	b := h.Hash("bz")
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}
