package slices

import "testing"

func TestIncludes(t *testing.T) {
	type args struct {
		slice []int
		item  int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "number1",
			args: args{
				slice: []int{2, 1, 3, 4, 5, 6},
				item:  4,
			},
			want: true,
		},
		{
			name: "number2",
			args: args{
				slice: []int{2, 1, 3, 4, 5, 6},
				item:  7,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Includes(tt.args.slice, tt.args.item); got != tt.want {
				t.Errorf("Includes() = %v, want %v", got, tt.want)
			}
		})
	}
}
