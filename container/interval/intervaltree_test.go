package interval

import (
	"reflect"
	"testing"
)

func eq[T any](t *testing.T, a, b []T) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("not eq %v|%v", a, b)
		t.FailNow()
	}
}

func TestQuery(t *testing.T) {
	interval1_3 := NewIntervalWithValue(1, 3, "1")
	interval2_3 := NewIntervalWithValue(2, 3, "2")
	interval3_3_3 := NewIntervalWithValue(3, 3, "3")
	interval3_3_4 := NewIntervalWithValue(3, 3, "4")
	interval4_6 := NewIntervalWithValue(4, 8, "5")
	interval5_7 := NewIntervalWithValue(5, 7, "6")
	interval6_9 := NewIntervalWithValue(6, 9, "7")
	tree := NewIntervalTree(interval1_3, interval2_3, interval3_3_3, interval3_3_4, interval4_6, interval5_7, interval6_9)
	tree.Print(10)
	eq(t, tree.Query(1), []string{interval1_3.value})
	eq(t, tree.Query(2), []string{interval2_3.value, interval1_3.value})
	eq(t, tree.Query(3), []string{interval3_3_4.value, interval2_3.value, interval1_3.value, interval3_3_3.value})
	eq(t, tree.Query(5), []string{interval5_7.value, interval4_6.value})
	eq(t, tree.Query(10), nil)
	eq(t, tree.Query(9), []string{interval6_9.value})
}
