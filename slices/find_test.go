package slices

import (
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	type args struct {
		slice     []int
		condition func(item int, index int, slice []int) bool
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name: "findNumber1",
			args: args{
				slice: []int{2, 3, 4, 13},
				condition: func(item int, _ int, _ []int) bool {
					return item == 14
				},
			},
			want:  0,
			want1: false,
		},
		{
			name: "findNumber2",
			args: args{
				slice: []int{2, 3, 4, 13},
				condition: func(item int, _ int, _ []int) bool {
					return item == 4
				},
			},
			want:  4,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Find(tt.args.slice, tt.args.condition)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
