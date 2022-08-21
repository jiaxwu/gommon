package slices

import (
	"fmt"
	"testing"
)

func TestForEach(t *testing.T) {
	type args struct {
		slice  []int
		action func(item int, index int, slice []int)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "number1",
			args: args{
				slice: []int{4, 5, 6},
				action: func(item int, index int, slice []int) {
					fmt.Printf("item: %d, index: %d\n", item, index)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ForEach(tt.args.slice, tt.args.action)
		})
	}
}
