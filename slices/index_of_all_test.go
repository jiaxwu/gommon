package slices

import (
	"reflect"
	"testing"
)

func TestIndexOfAll(t *testing.T) {
	type args struct {
		slice []int
		item  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "number1",
			args: args{
				slice: []int{3, 4, 5, 3, 2},
				item:  4,
			},
			want: []int{1},
		},
		{
			name: "number2",
			args: args{
				slice: []int{3, 4, 5, 3, 2},
				item:  3,
			},
			want: []int{0, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IndexOfAll(tt.args.slice, tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
