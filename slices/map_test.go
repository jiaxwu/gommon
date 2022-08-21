package slices

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	type args struct {
		slice  []int
		mapper func(item int, index int, slice []int) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "intToString",
			args: args{
				slice: []int{2, 3, 4, 13},
				mapper: func(item int, index int, slice []int) string {
					return fmt.Sprintf("%d:%d", item, index)
				},
			},
			want: []string{"2:0", "3:1", "4:2", "13:3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.slice, tt.args.mapper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
