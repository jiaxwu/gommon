package tinylfu

import (
	"testing"

	"github.com/jiaxwu/gommon/cache"
)

func TestCache_Put(t *testing.T) {
	c := New[string, int](func(key string) []byte {
		return []byte(key)
	}, 3)
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Get("11")
	c.Put("44", 8)

	value, ok := c.Get("33")
	if value != 0 || ok {
		t.Errorf("Get() = %v, want %v", ok, false)
	}
}

func TestCache_OnEvict(t *testing.T) {
	c := New[string, int](func(key string) []byte {
		return []byte(key)
	}, 3)
	c.SetOnEvict(func(entry *cache.Entry[string, int]) {
		if entry.Key != "22" || entry.Value != 6 {
			t.Errorf("OnEvict() = %v, want %v", entry.Key, "22")
		}
	})
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Get("11")
	c.Put("44", 8)

	value, ok := c.Get("33")
	if value != 0 || ok {
		t.Errorf("Get() = %v, want %v", ok, false)
	}
}

func TestCache_Clear(t *testing.T) {
	c := New[string, int](func(key string) []byte {
		return []byte(key)
	}, 3)
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Get("11")
	c.Put("44", 8)

	value, ok := c.Get("33")
	if value != 0 || ok {
		t.Errorf("Get() = %v, want %v", ok, false)
	}

	value, ok = c.Get("11")
	if value != 5 || !ok {
		t.Errorf("Get() = %v, want %v", ok, true)
	}

	c.Clear(false)
	value, ok = c.Get("11")
	if value != 0 || ok {
		t.Errorf("Get() = %v, want %v", ok, false)
	}
}

func TestCache_Peek(t *testing.T) {
	c := New[string, int](func(key string) []byte {
		return []byte(key)
	}, 3)
	c.Put("11", 5)
	c.Put("22", 6)
	c.Put("33", 7)
	c.Peek("11")
	c.Put("44", 8)

	value, ok := c.Get("33")
	if value != 0 || ok {
		t.Errorf("Get() = %v, want %v", ok, false)
	}
}

func TestCache_Remove(t *testing.T) {
	c := New[string, int](func(key string) []byte {
		return []byte(key)
	}, 3)
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

func TestCache_Entries(t *testing.T) {
	c := New[string, int](func(key string) []byte {
		return []byte(key)
	}, 3)
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

// hitRate=0.371382
// Fuzz基本随机，对缓存测试并不科学
func FuzzHitRate(f *testing.F) {
	seeds := []string{"abc", "bbb", "0", "1", "", "zdas", "xzasd", "1312", "0", "0", "0", "0", "1", "1", "1"}
	for _, seed := range seeds {
		f.Add(seed)
	}
	n := 100000
	mul := 50
	c := New[string, int](func(key string) []byte {
		return []byte(key)
	}, n)
	count := 0
	hitCount := 0
	m := map[string]int{}
	f.Fuzz(func(t *testing.T, a string) {
		count++
		m[a]++
		_, exists := c.Get(a)
		if exists {
			hitCount++
		} else {
			c.Put(a, 0)
		}
		if count == n*mul {
			hitRate := float64(hitCount) / float64(count)
			t.Errorf("count=%v, hitCount=%v, hitRate=%f, items=%v", count, hitCount, hitRate, len(m))
			t.SkipNow()
		}
	})
}
