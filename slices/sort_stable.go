package slices

import "sort"

// SortStable sorts the slice given the provided condition function,
// keeping equal elements in their original order.
func SortStable[T any](slice []T, condition func(item1, item2 T) bool) []T {
	sort.SliceStable(slice, func(i, j int) bool {
		return condition(slice[i], slice[j])
	})
	return slice
}
