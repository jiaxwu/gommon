package round

import (
	"math"
	"testing"
)

// uint overflow
func TestUintOverflow(t *testing.T) {
	var in uint = math.MaxUint64
	var out uint = math.MaxUint64 - 1
	if in-out != 1 {
		t.Errorf("want %d, but %d", 1, in-out)
	}
	in++
	if in-out != 2 {
		t.Errorf("want %d, but %d", 2, in-out)
	}
	out++
	if in-out != 1 {
		t.Errorf("want %d, but %d", 1, in-out)
	}
}
