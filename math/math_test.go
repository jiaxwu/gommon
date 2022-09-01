package math

import (
	"math"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name string
		n    uint
		m    uint
		per  uint
		last uint
	}{
		{
			name: "1",
			n:    10,
			m:    3,
			per:  4,
			last: 2,
		},
		{
			name: "2",
			n:    10,
			m:    4,
			per:  3,
			last: 1,
		},
		{
			name: "3",
			n:    11,
			m:    3,
			per:  4,
			last: 3,
		},
		{
			name: "4",
			n:    12,
			m:    3,
			per:  4,
			last: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			per, last := Split(tt.n, tt.m)
			if per != tt.per {
				t.Errorf("Split() per = %v, expected %v", per, tt.per)
			}
			if last != tt.last {
				t.Errorf("Split() last = %v, expected %v", last, tt.last)
			}
		})
	}
}

func TestIsPowOf2(t *testing.T) {
	tests := []struct {
		name string
		x    uint64
		is   bool
	}{
		{
			name: "0",
			x:    0,
			is:   false,
		},
		{
			name: "1",
			x:    1,
			is:   true,
		},
		{
			name: "2",
			x:    2,
			is:   true,
		},
		{
			name: "4",
			x:    4,
			is:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := IsPowOf2(tt.x)
			if is != tt.is {
				t.Errorf("IsPowOf2() is = %v, expected %v", is, tt.is)
			}
		})
	}
}

func TestRoundUpPowOf2(t *testing.T) {
	tests := []struct {
		name string
		x    uint64
		up   uint64
	}{
		{
			name: "0",
			x:    0,
			up:   0,
		},
		{
			name: "1",
			x:    1,
			up:   1,
		},
		{
			name: "2",
			x:    3,
			up:   4,
		},
		{
			name: "4",
			x:    10,
			up:   16,
		},
		{
			name: "math.MaxUint64",
			x:    math.MaxUint64,
			up:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			up := RoundUpPowOf2(tt.x)
			if up != tt.up {
				t.Errorf("RoundUpPowOf2() up = %v, expected %v", up, tt.up)
			}
		})
	}
}

func TestRoundDownPowOf2(t *testing.T) {
	tests := []struct {
		name string
		x    uint64
		down uint64
	}{
		{
			name: "1",
			x:    1,
			down: 1,
		},
		{
			name: "2",
			x:    3,
			down: 2,
		},
		{
			name: "4",
			x:    10,
			down: 8,
		},
		{
			name: "math.MaxUint64",
			x:    math.MaxUint64,
			down: 1 << 63,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			down := RoundDownPowOf2(tt.x)
			if down != tt.down {
				t.Errorf("RoundDownPowOf2() down = %v, expected %v", down, tt.down)
			}
		})
	}
}

func TestFindLastBitSet(t *testing.T) {
	type entry = struct {
		nlz, ntz, pop int
	}

	// tab contains results for all uint8 values
	var tab [256]entry

	tab[0] = entry{8, 8, 0}
	for i := 1; i < len(tab); i++ {
		// nlz
		x := i // x != 0
		n := 0
		for x&0x80 == 0 {
			n++
			x <<= 1
		}
		tab[i].nlz = n

		// ntz
		x = i // x != 0
		n = 0
		for x&1 == 0 {
			n++
			x >>= 1
		}
		tab[i].ntz = n

		// pop
		x = i // x != 0
		n = 0
		for x != 0 {
			n += int(x & 1)
			x >>= 1
		}
		tab[i].pop = n
	}
	for i := 0; i < 256; i++ {
		len := 8 - tab[i].nlz
		for k := 0; k < 64-8; k++ {
			x := uint64(i) << uint(k)
			want := 0
			if x != 0 {
				want = len + k
			}
			if x <= 1<<8-1 {
				got := FindLastBitSet(uint8(x))
				if got != want {
					t.Fatalf("Len8(%#02x) == %d; want %d", x, got, want)
				}
			}

			if x <= 1<<16-1 {
				got := FindLastBitSet(uint16(x))
				if got != want {
					t.Fatalf("Len16(%#04x) == %d; want %d", x, got, want)
				}
			}

			if x <= 1<<32-1 {
				got := FindLastBitSet(uint32(x))
				if got != want {
					t.Fatalf("Len32(%#08x) == %d; want %d", x, got, want)
				}
			}

			if x <= 1<<64-1 {
				got := FindLastBitSet(uint64(x))
				if got != want {
					t.Fatalf("Len64(%#016x) == %d; want %d", x, got, want)
				}
				got = FindLastBitSet(uint(x))
				if got != want {
					t.Fatalf("Len(%#016x) == %d; want %d", x, got, want)
				}
			}
		}
	}
}
