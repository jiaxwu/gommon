package slices

import "testing"

func TestEqual(t *testing.T) {
	type args struct {
		slice1 []int
		slice2 []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "number1",
			args: args{
				slice1: []int{1, 2, 3},
				slice2: []int{1, 2, 3},
			},
			want: true,
		},
		{
			name: "number2",
			args: args{
				slice1: []int{1, 2, 3},
				slice2: []int{2, 2, 3},
			},
			want: false,
		},
		{
			name: "number3",
			args: args{
				slice1: []int{1, 2, 3},
				slice2: []int{1, 2, 3, 4},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.slice1, tt.args.slice2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualFunc(t *testing.T) {
	type args struct {
		slice1 []int
		slice2 []int
		f      func(item1, item2 int) bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "number1",
			args: args{
				slice1: []int{1, 2, 3},
				slice2: []int{1, 2, 3},
				f: func(item1, item2 int) bool {
					return item1 == item2
				},
			},
			want: true,
		},
		{
			name: "number2",
			args: args{
				slice1: []int{1, 2, 3},
				slice2: []int{2, 2, 3},
				f: func(item1, item2 int) bool {
					return item1 == item2
				},
			},
			want: false,
		},
		{
			name: "number3",
			args: args{
				slice1: []int{1, 2, 3},
				slice2: []int{1, 2, 3, 4},
				f: func(item1, item2 int) bool {
					return item1 == item2
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualFunc(tt.args.slice1, tt.args.slice2, tt.args.f); got != tt.want {
				t.Errorf("EqualFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
