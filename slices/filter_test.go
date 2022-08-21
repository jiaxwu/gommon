package slices

import (
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	type args struct {
		slice  []int
		filter func(item int, index int, slice []int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "number1",
			args: args{
				slice: []int{18, 19, 6, 3, 43, 1, 32},
				filter: func(item int, index int, _ []int) bool {
					return item > 18 && index > 2
				},
			},
			want: []int{43, 32},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.slice, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
