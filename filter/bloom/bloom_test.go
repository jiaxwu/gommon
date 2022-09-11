package bloom

import (
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	f := New(10, 0.01)
	f.AddString("10")
	f.AddString("44")
	f.AddString("66")
	if !f.ContainsString("10") {
		t.Errorf("want %v, but %v", true, f.ContainsString("10"))
	}
	if !f.ContainsString("44") {
		t.Errorf("want %v, but %v", true, f.ContainsString("10"))
	}
	if !f.ContainsString("66") {
		t.Errorf("want %v, but %v", true, f.ContainsString("10"))
	}
	if f.ContainsString("55") {
		t.Errorf("want %v, but %v", false, f.ContainsString("10"))
	}
}

func BenchmarkAddAndContains(b *testing.B) {
	buf := make([]byte, 8192)
	for length := 1; length <= cap(buf); length *= 2 {
		b.Run(strconv.Itoa(length), func(b *testing.B) {
			f := New(uint64(b.N), 0.0001)
			buf = buf[:length]
			b.SetBytes(int64(length))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f.Add(buf)
				f.Contains(buf)
			}
		})
	}
}
