package slices

import "testing"

func TestEvery(t *testing.T) {
	type args struct {
		slice     []int
		condition func(item int, index int, slice []int) bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "number1",
			args: args{
				slice: []int{2, 2, 2, 2},
				condition: func(item int, _ int, _ []int) bool {
					return item == 2
				},
			},
			want: true,
		},
		{
			name: "number2",
			args: args{
				slice: []int{2, 3, 4, 13},
				condition: func(item int, _ int, _ []int) bool {
					return item == 2
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Every(tt.args.slice, tt.args.condition); got != tt.want {
				t.Errorf("Every() = %v, want %v", got, tt.want)
			}
		})
	}
}
