package slices

import (
	"reflect"
	"strconv"
	"testing"
)

func TestReduce(t *testing.T) {
	type args struct {
		slice  []int
		reduce func(total int, item int, index int, slice []int) int
		init   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "number1",
			args: args{
				slice: []int{3, 4, 5, 3},
				reduce: func(total int, item int, index int, slice []int) int {
					return total + item
				},
				init: 0,
			},
			want: 15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.slice, tt.args.reduce, tt.args.init); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduce2(t *testing.T) {
	type args struct {
		slice  []int
		reduce func(total string, item int, index int, slice []int) string
		init   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "numberToString",
			args: args{
				slice: []int{1, 2, 3, 4, 5},
				reduce: func(total string, item int, index int, slice []int) string {
					return total + strconv.Itoa(item)
				},
				init: "",
			},
			want: "12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.slice, tt.args.reduce, tt.args.init); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}
