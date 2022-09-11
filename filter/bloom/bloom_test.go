package bloom

import (
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
