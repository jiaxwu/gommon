package slices

import "testing"

func TestLastIndexOf(t *testing.T) {
	type args struct {
		slice []int
		item  int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "number1",
			args: args{
				slice: []int{3, 4, 5, 3, 2},
				item:  3,
			},
			want: 3,
		},
		{
			name: "number4",
			args: args{
				slice: []int{3, 4, 5, 3, 2},
				item:  4,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LastIndexOf(tt.args.slice, tt.args.item); got != tt.want {
				t.Errorf("LastIndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
