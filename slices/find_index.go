package slices

// FindIndex first index that meet the condition
func FindIndex[T any](slice []T, condition func(item T, index int, slice []T) bool) int {
	for index, item := range slice {
		if condition(item, index, slice) {
			return index
		}
	}
	return -1
}
