package slices

import (
	"reflect"
	"testing"
)

func TestUnshift(t *testing.T) {
	type args struct {
		slice []int
		items []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "number",
			args: args{
				slice: []int{3, 4, 5},
				items: []int{1, 2},
			},
			want: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unshift(tt.args.slice, tt.args.items...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unshift() = %v, want %v", got, tt.want)
			}
		})
	}
}
