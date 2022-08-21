package slices

import (
	"reflect"
	"testing"
)

func TestFrequencies(t *testing.T) {
	type args struct {
		slice []int
	}
	tests := []struct {
		name string
		args args
		want map[int]int
	}{
		{
			name: "number1",
			args: args{
				slice: []int{1, 2, 3, 4, 3, 2, 1, 5},
			},
			want: map[int]int{
				1: 2,
				2: 2,
				3: 2,
				4: 1,
				5: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Frequencies(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Frequencies() = %v, want %v", got, tt.want)
			}
		})
	}
}
