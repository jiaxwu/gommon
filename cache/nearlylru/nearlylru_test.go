package nearlylru

import (
	"os"
	"strings"
	"testing"
)

// nearlylru_test.go:58: samples=5, cachePercentage=0.1%, count=206048, hitCount=26545, hitRate=12.88%
// nearlylru_test.go:58: samples=5, cachePercentage=0.3%, count=206048, hitCount=56550, hitRate=27.45%
// nearlylru_test.go:58: samples=5, cachePercentage=0.5%, count=206048, hitCount=84843, hitRate=41.18%
// nearlylru_test.go:58: samples=5, cachePercentage=0.7%, count=206048, hitCount=108567, hitRate=52.69%
// nearlylru_test.go:58: samples=5, cachePercentage=1.0%, count=206048, hitCount=139522, hitRate=67.71%
// nearlylru_test.go:58: samples=5, cachePercentage=2.0%, count=206048, hitCount=182253, hitRate=88.45%
// nearlylru_test.go:58: samples=5, cachePercentage=3.0%, count=206048, hitCount=189181, hitRate=91.81%
// nearlylru_test.go:58: samples=5, cachePercentage=5.0%, count=206048, hitCount=192501, hitRate=93.43%
// nearlylru_test.go:58: samples=5, cachePercentage=10.0%, count=206048, hitCount=192842, hitRate=93.59%
// nearlylru_test.go:58: samples=10, cachePercentage=0.1%, count=206048, hitCount=26568, hitRate=12.89%
// nearlylru_test.go:58: samples=10, cachePercentage=0.3%, count=206048, hitCount=57342, hitRate=27.83%
// nearlylru_test.go:58: samples=10, cachePercentage=0.5%, count=206048, hitCount=85452, hitRate=41.47%
// nearlylru_test.go:58: samples=10, cachePercentage=0.7%, count=206048, hitCount=111486, hitRate=54.11%
// nearlylru_test.go:58: samples=10, cachePercentage=1.0%, count=206048, hitCount=143822, hitRate=69.80%
// nearlylru_test.go:58: samples=10, cachePercentage=2.0%, count=206048, hitCount=184845, hitRate=89.71%
// nearlylru_test.go:58: samples=10, cachePercentage=3.0%, count=206048, hitCount=190227, hitRate=92.32%
// nearlylru_test.go:58: samples=10, cachePercentage=5.0%, count=206048, hitCount=192551, hitRate=93.45%
// nearlylru_test.go:58: samples=10, cachePercentage=10.0%, count=206048, hitCount=192842, hitRate=93.59%
// nearlylru_test.go:58: samples=20, cachePercentage=0.1%, count=206048, hitCount=26548, hitRate=12.88%
// nearlylru_test.go:58: samples=20, cachePercentage=0.3%, count=206048, hitCount=57196, hitRate=27.76%
// nearlylru_test.go:58: samples=20, cachePercentage=0.5%, count=206048, hitCount=86258, hitRate=41.86%
// nearlylru_test.go:58: samples=20, cachePercentage=0.7%, count=206048, hitCount=112824, hitRate=54.76%
// nearlylru_test.go:58: samples=20, cachePercentage=1.0%, count=206048, hitCount=146299, hitRate=71.00%
// nearlylru_test.go:58: samples=20, cachePercentage=2.0%, count=206048, hitCount=186292, hitRate=90.41%
// nearlylru_test.go:58: samples=20, cachePercentage=3.0%, count=206048, hitCount=190549, hitRate=92.48%
// nearlylru_test.go:58: samples=20, cachePercentage=5.0%, count=206048, hitCount=192597, hitRate=93.47%
// nearlylru_test.go:58: samples=20, cachePercentage=10.0%, count=206048, hitCount=192842, hitRate=93.59%
// nearlylru_test.go:58: samples=50, cachePercentage=0.1%, count=206048, hitCount=26678, hitRate=12.95%
// nearlylru_test.go:58: samples=50, cachePercentage=0.3%, count=206048, hitCount=57943, hitRate=28.12%
// nearlylru_test.go:58: samples=50, cachePercentage=0.5%, count=206048, hitCount=87033, hitRate=42.24%
// nearlylru_test.go:58: samples=50, cachePercentage=0.7%, count=206048, hitCount=113703, hitRate=55.18%
// nearlylru_test.go:58: samples=50, cachePercentage=1.0%, count=206048, hitCount=147607, hitRate=71.64%
// nearlylru_test.go:58: samples=50, cachePercentage=2.0%, count=206048, hitCount=186958, hitRate=90.74%
// nearlylru_test.go:58: samples=50, cachePercentage=3.0%, count=206048, hitCount=190593, hitRate=92.50%
// nearlylru_test.go:58: samples=50, cachePercentage=5.0%, count=206048, hitCount=192595, hitRate=93.47%
// nearlylru_test.go:58: samples=50, cachePercentage=10.0%, count=206048, hitCount=192842, hitRate=93.59%
func TestHitRate(t *testing.T) {
	dataset, err := os.ReadFile("../dataset")
	if err != nil {
		t.Errorf("read dataset error %v", err)
	}
	reqs := strings.Split(string(dataset), ",")
	testHitRates(t, reqs, 5)
	testHitRates(t, reqs, 10)
	testHitRates(t, reqs, 20)
	testHitRates(t, reqs, 50)
}

func testHitRates(t *testing.T, reqs []string, samples int) {
	testHitRate(t, reqs, 0.001, samples)
	testHitRate(t, reqs, 0.003, samples)
	testHitRate(t, reqs, 0.005, samples)
	testHitRate(t, reqs, 0.007, samples)
	testHitRate(t, reqs, 0.01, samples)
	testHitRate(t, reqs, 0.02, samples)
	testHitRate(t, reqs, 0.03, samples)
	testHitRate(t, reqs, 0.05, samples)
	testHitRate(t, reqs, 0.1, samples)
}

func testHitRate(t *testing.T, reqs []string, cachePercentage float64, samples int) {
	count := len(reqs)
	n := int(float64(count) * cachePercentage)
	c := New[string, int](n)
	c.SetSamples(samples)
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
	t.Logf("samples=%v, cachePercentage=%.1f%%, count=%v, hitCount=%v, hitRate=%.2f%%", samples, cachePercentage*100, count, hitCount, hitRate*100)
	// t.Logf("|%v|%.1f%%| %v|%.2f%%|", samples, cachePercentage*100, hitCount, hitRate*100)
}
