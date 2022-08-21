package slices

import (
	"reflect"
	"testing"
)

func TestFill(t *testing.T) {
	type args struct {
		slice    []int
		itemFunc func(item int, index int, slice []int) int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "number1",
			args: args{
				slice: make([]int, 10),
				itemFunc: func(item int, index int, slice []int) int {
					return index * index
				},
			},
			want: []int{0, 1, 4, 9, 16, 25, 36, 49, 64, 81},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fill(tt.args.slice, tt.args.itemFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fill() = %v, want %v", got, tt.want)
			}
		})
	}
}
