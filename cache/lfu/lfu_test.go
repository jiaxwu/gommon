package lfu

import (
	"os"
	"strings"
	"testing"

	"github.com/jiaxwu/gommon/cache"
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
	c.SetOnEvict(func(entry *cache.Entry[string, int]) {
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

// lfu_test.go:172: cachePercentage=0.1%, count=206048, hitCount=28322, hitRate=13.75%
// lfu_test.go:172: cachePercentage=0.3%, count=206048, hitCount=59827, hitRate=29.04%
// lfu_test.go:172: cachePercentage=0.5%, count=206048, hitCount=88984, hitRate=43.19%
// lfu_test.go:172: cachePercentage=0.7%, count=206048, hitCount=115660, hitRate=56.13%
// lfu_test.go:172: cachePercentage=1.0%, count=206048, hitCount=149970, hitRate=72.78%
// lfu_test.go:172: cachePercentage=2.0%, count=206048, hitCount=187426, hitRate=90.96%
// lfu_test.go:172: cachePercentage=3.0%, count=206048, hitCount=190666, hitRate=92.53%
// lfu_test.go:172: cachePercentage=5.0%, count=206048, hitCount=192569, hitRate=93.46%
// lfu_test.go:172: cachePercentage=10.0%, count=206048, hitCount=192842, hitRate=93.59%
func TestHitRate(t *testing.T) {
	dataset, err := os.ReadFile("../dataset")
	if err != nil {
		t.Errorf("read dataset error %v", err)
	}
	reqs := strings.Split(string(dataset), ",")
	testHitRate(t, reqs, 0.001)
	testHitRate(t, reqs, 0.003)
	testHitRate(t, reqs, 0.005)
	testHitRate(t, reqs, 0.007)
	testHitRate(t, reqs, 0.01)
	testHitRate(t, reqs, 0.02)
	testHitRate(t, reqs, 0.03)
	testHitRate(t, reqs, 0.05)
	testHitRate(t, reqs, 0.1)
}

func testHitRate(t *testing.T, reqs []string, cachePercentage float64) {
	count := len(reqs)
	n := int(float64(count) * cachePercentage)
	c := New[string, int](n)
	hitCount := 0
	for _, req := range reqs {
		_, exists := c.Get(req)
		if exists {
			hitCount++
		} else {
			c.Put(req, 0)
		}
	}
	hitRate := float64(hitCount) / float64(count)
	t.Logf("cachePercentage=%.1f%%, count=%v, hitCount=%v, hitRate=%.2f%%", cachePercentage*100, count, hitCount, hitRate*100)
}
