package slices

import "sort"

// Sort sorts the slice given the provided condition function.
func Sort[T any](slice []T, condition func(item1, item2 T) bool) []T {
	sort.Slice(slice, func(i, j int) bool {
		return condition(slice[i], slice[j])
	})
	return slice
}
