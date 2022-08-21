package slices

import (
	"reflect"
	"testing"
)

type Item struct {
	Index int
	Value int
}

func TestSortStable(t *testing.T) {
	type args struct {
		slice     []Item
		condition func(item1, item2 Item) bool
	}
	tests := []struct {
		name string
		args args
		want []Item
	}{
		{
			name: "item1",
			args: args{
				slice: []Item{{1, 1}, {2, 3}, {3, 5}, {4, 3},
					{5, 4}, {6, 1}, {7, 1}, {8, 1}},
				condition: func(item1, item2 Item) bool {
					return item1.Value < item2.Value
				},
			},
			want: []Item{{1, 1}, {6, 1}, {7, 1}, {8, 1},
				{2, 3}, {4, 3}, {5, 4}, {3, 5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortStable(tt.args.slice, tt.args.condition); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortStable() = %v, want %v", got, tt.want)
			}
		})
	}
}
