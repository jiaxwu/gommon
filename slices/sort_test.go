package slices

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	type args struct {
		slice     []int
		condition func(item1, item2 int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "number1",
			args: args{
				slice: []int{3, 4, 2, 1, 5, 7, 6, 8},
				condition: func(item1, item2 int) bool {
					return item1 < item2
				},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sort(tt.args.slice, tt.args.condition); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}
