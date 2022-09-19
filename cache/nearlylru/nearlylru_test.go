package nearlylru

import "testing"

// hitRate=0.388525
func FuzzHitRate(f *testing.F) {
	seeds := []string{"abc", "bbb", "0", "1", ""}
	for _, seed := range seeds {
		f.Add(seed)
	}
	n := 100000
	mul := 20
	c := New[string, int](n)
	count := 0
	hitCount := 0
	f.Fuzz(func(t *testing.T, a string) {
		count++
		_, exists := c.Get(a)
		if exists {
			hitCount++
		} else {
			c.Put(a, 0)
		}
		if count == n*mul {
			hitRate := float64(hitCount) / float64(count)
			t.Errorf("count=%v, hitCount=%v, hitRate=%f", count, hitCount, hitRate)
			t.SkipNow()
		}
	})
}
