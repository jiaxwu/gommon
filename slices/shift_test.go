package slices

import (
	"reflect"
	"testing"
)

func TestShift(t *testing.T) {
	type args struct {
		slice []int
	}
	tests := []struct {
		name  string
		args  args
		want  []int
		want1 int
	}{
		{
			name: "number1",
			args: args{
				slice: []int{3, 4, 5, 6},
			},
			want:  []int{4, 5, 6},
			want1: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Shift(tt.args.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shift() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Shift() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
