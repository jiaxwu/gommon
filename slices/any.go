package slices

// Any item is meet the condition
func Any[T any](slice []T, condition func(item T, index int, slice []T) bool) bool {
	for index, item := range slice {
		if condition(item, index, slice) {
			return true
		}
	}
	return false
}
