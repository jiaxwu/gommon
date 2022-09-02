package round

import (
	"testing"
)

var (
	Len  = 8
	Mask = Len - 1
	In   = 8 - 5
)

// % len
func BenchmarkModLen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = In % Len
	}
}

// & Mask
func BenchmarkAndMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = In & Mask
	}
}
