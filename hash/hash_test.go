package hash

import (
	"strconv"
	"testing"

	"github.com/jiaxwu/gommon/conv"
)

type Key struct {
	S string
	B int
}

func TestStruct(t *testing.T) {
	h := New()
	a := h.Sum64(conv.MustMarshal(Key{
		S: "ab",
		B: 6,
	}))
	b := h.Sum64(conv.MustMarshal(Key{
		S: "ab",
		B: 6,
	}))
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func TestString(t *testing.T) {
	h := New()
	a := h.Sum64String("bz")
	b := h.Sum64String("bz")
	if a != b {
		t.Errorf("want %v, but %v", a, b)
	}
}

func Benchmark64(b *testing.B) {
	buf := make([]byte, 8192)
	for length := 1; length <= cap(buf); length *= 2 {
		b.Run(strconv.Itoa(length), func(b *testing.B) {
			h := New()
			buf = buf[:length]
			b.SetBytes(int64(length))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Sum64(buf)
			}
		})
	}
}
