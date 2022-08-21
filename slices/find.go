package slices

// Find first item that meet the condition
func Find[T any](slice []T, condition func(item T, index int, slice []T) bool) (T, bool) {
	for index, item := range slice {
		if condition(item, index, slice) {
			return item, true
		}
	}
	var t T
	return t, false
}
