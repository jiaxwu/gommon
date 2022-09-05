package lfu

import (
	"testing"
)

func TestCache_Put(t *testing.T) {
	c := New[string, int](3)
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Get("11")
	c.Get("33")
	c.Put("44", 8)

	value, ok := c.Get("22")
	if value != 0 || ok {
		t.Errorf("Put() = %v, want %v", ok, false)
	}
}

func TestCache_OnEvict(t *testing.T) {
	c := New[string, int](3)
	c.SetOnEvict(func(entry *Entry[string, int]) {
		if entry.Key != "22" || entry.Value != 6 {
			t.Errorf("OnEvict() = %v, want %v", entry.Key, "22")
		}
	})
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Get("11")
	c.Get("33")
	c.Put("44", 8)

	value, ok := c.Get("22")
	if value != 0 || ok {
		t.Errorf("Put() = %v, want %v", ok, false)
	}
}

func TestCache_Clear(t *testing.T) {
	c := New[string, int](3)
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Get("11")
	c.Get("33")
	c.Put("44", 8)

	value, ok := c.Get("22")
	if value != 0 || ok {
		t.Errorf("Put() = %v, want %v", ok, false)
	}

	value, ok = c.Get("11")
	if value != 5 || !ok {
		t.Errorf("Put() = %v, want %v", ok, true)
	}

	c.Clear(false)
	value, ok = c.Get("11")
	if value != 0 || ok {
		t.Errorf("Put() = %v, want %v", ok, false)
	}
}

func TestCache_Peek(t *testing.T) {
	c := New[string, int](3)
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Peek("11")
	c.Peek("33")
	c.Put("44", 8)

	value, ok := c.Get("33")
	if value != 0 || ok {
		t.Errorf("Put() = %v, want %v", ok, false)
	}
}

func TestCache_Remove(t *testing.T) {
	c := New[string, int](3)
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Remove("22")
	c.Put("44", 8)

	value, ok := c.Get("22")
	if value != 0 || ok {
		t.Errorf("Get() = %v, want %v", value, 0)
	}
}

func TestCache_Evict(t *testing.T) {
	c := New[string, int](3)
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Evict()
	c.Put("44", 8)

	value, ok := c.Get("33")
	if value != 0 || ok {
		t.Errorf("Get() = %v, want %v", value, 0)
	}
}

func TestCache_Entries(t *testing.T) {
	c := New[string, int](3)
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Get("11")
	c.Put("44", 8)

	entries := c.Entries()
	keys := []string{"44", "22", "11"}
	for i, entry := range entries {
		if entry.Key != keys[i] {
			t.Errorf("Get() = %v, want %v", entry.Key, keys[i])
		}
	}
}