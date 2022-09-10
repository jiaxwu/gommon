package cm

import (
	"testing"
)

func TestCount(t *testing.T) {
	cm := New(10)
	cm.Add(10)
	cm.Add(51151)
	cm.Add(321)
	cm.Add(10)
	cm.Add(10)
	cm.Add(321)
	if cm.Estimate(10) != 3 {
		t.Errorf("want %v, but %d", 3, cm.Estimate(10))
	}
	if cm.Estimate(321) != 2 {
		t.Errorf("want %v, but %d", 2, cm.Estimate(321))
	}
	if cm.Estimate(51151) != 1 {
		t.Errorf("want %v, but %d", 1, cm.Estimate(1))
	}

	for i := 0; i < 100; i++ {
		cm.Add(10)
	}
	if cm.Estimate(10) != 15 {
		t.Errorf("want %v, but %d", 15, cm.Estimate(10))
	}
	for i := 0; i < 100; i++ {
		cm.Add(5)
	}
	if cm.Estimate(10) != 7 {
		t.Errorf("want %v, but %d", 7, cm.Estimate(10))
	}
	for i := 0; i < 100; i++ {
		cm.Add(uint64(1))
	}
	if cm.Estimate(10) != 7 {
		t.Errorf("want %v, but %d", 7, cm.Estimate(10))
	}
}
