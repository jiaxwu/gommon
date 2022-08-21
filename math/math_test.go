package math

import "testing"

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
