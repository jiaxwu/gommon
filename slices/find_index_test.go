package slices

import "testing"

func TestFindIndex(t *testing.T) {
	type args struct {
		slice     []int
		condition func(item int, index int, slice []int) bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "findNumber1",
			args: args{
				slice: []int{2, 3, 4, 13},
				condition: func(item int, _ int, _ []int) bool {
					return item == 14
				},
			},
			want: -1,
		},
		{
			name: "findNumber2",
			args: args{
				slice: []int{2, 3, 4, 13},
				condition: func(item int, _ int, _ []int) bool {
					return item == 4
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindIndex(tt.args.slice, tt.args.condition); got != tt.want {
				t.Errorf("FindIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
