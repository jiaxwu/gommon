package slices

import (
	"testing"
)

func TestShuffle(t *testing.T) {
	type args struct {
		slice []int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "number1",
			args: args{
				slice: []int{1, 2, 3, 4, 5, 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Shuffle(tt.args.slice)
		})
	}
}
