package slices

import "testing"

func TestFirstNotZero(t *testing.T) {
	type args struct {
		slice []int
	}
	type want struct {
		v  int
		ok bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "number1",
			args: args{
				slice: []int{0, 2, 3, 4, 13},
			},
			want: want{
				v:  2,
				ok: true,
			},
		},
		{
			name: "number2",
			args: args{
				slice: []int{4, 3, 4, 13},
			},
			want: want{
				v:  4,
				ok: true,
			},
		},
		{
			name: "number3",
			args: args{
				slice: []int{0, 0, 0, 0},
			},
			want: want{
				ok: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := FirstNotZero(tt.args.slice...); got != tt.want.v && ok != tt.want.ok {
				t.Errorf("FirstNotZero() = %v, %v, want %v and %v", got, ok, tt.want.v, tt.want.ok)
			}
		})
	}
}
