package slices

// FindLastIndex that meet the condition
func FindLastIndex[T any](slice []T, condition func(item T, index int, slice []T) bool) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if condition(slice[i], i, slice) {
			return i
		}
	}
	return -1
}
