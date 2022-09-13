package cm

import "testing"

func TestCount(t *testing.T) {
	cm := New[uint8](1000, 10, 0.001)
	cm.IncString("10")
	cm.IncString("51151")
	cm.IncString("321")
	cm.IncString("10")
	cm.IncString("10")
	cm.IncString("321")
	if cm.EstimateString("10") != 3 {
		t.Errorf("want %v, but %d", 3, cm.EstimateString("10"))
	}
	if cm.EstimateString("321") != 2 {
		t.Errorf("want %v, but %d", 2, cm.EstimateString("321"))
	}
	if cm.EstimateString("51151") != 1 {
		t.Errorf("want %v, but %d", 1, cm.EstimateString("1"))
	}

	cm.AddString("10", 100)
	if cm.EstimateString("10") != 103 {
		t.Errorf("want %v, but %d", 103, cm.EstimateString("10"))
	}
	cm.AddString("10", 254)
	if cm.EstimateString("10") != 255 {
		t.Errorf("want %v, but %d", 255, cm.EstimateString("10"))
	}
	cm.AddString("5", 100)
	if cm.EstimateString("5") != 100 {
		t.Errorf("want %v, but %d", 100, cm.EstimateString("5"))
	}
	cm.AddString("1", 100)
	if cm.EstimateString("1") != 100 {
		t.Errorf("want %v, but %d", 100, cm.EstimateString("1"))
	}

	cm.Attenuation(2)
	if cm.EstimateString("10") != 127 {
		t.Errorf("want %v, but %d", 127, cm.EstimateString("10"))
	}
	if cm.EstimateString("5") != 50 {
		t.Errorf("want %v, but %d", 50, cm.EstimateString("5"))
	}
	if cm.EstimateString("1") != 50 {
		t.Errorf("want %v, but %d", 50, cm.EstimateString("1"))
	}
}
