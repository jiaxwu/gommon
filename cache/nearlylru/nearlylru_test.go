package nearlylru

import (
	"os"
	"strings"
	"testing"
)

// nearlylru_test.go:49: cachePercentage=0.1%, count=206048, hitCount=26610, hitRate=12.91%
// nearlylru_test.go:49: cachePercentage=0.3%, count=206048, hitCount=56500, hitRate=27.42%
// nearlylru_test.go:49: cachePercentage=0.5%, count=206048, hitCount=83817, hitRate=40.68%
// nearlylru_test.go:49: cachePercentage=0.7%, count=206048, hitCount=107777, hitRate=52.31%
// nearlylru_test.go:49: cachePercentage=1.0%, count=206048, hitCount=139395, hitRate=67.65%
// nearlylru_test.go:49: cachePercentage=2.0%, count=206048, hitCount=181211, hitRate=87.95%
// nearlylru_test.go:49: cachePercentage=3.0%, count=206048, hitCount=189003, hitRate=91.73%
// nearlylru_test.go:49: cachePercentage=5.0%, count=206048, hitCount=192485, hitRate=93.42%
// nearlylru_test.go:49: cachePercentage=10.0%, count=206048, hitCount=192842, hitRate=93.59%
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
