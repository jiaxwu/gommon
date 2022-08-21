package slices

// Every item is meet the condition
func Every[T any](slice []T, condition func(item T, index int, slice []T) bool) bool {
	for index, item := range slice {
		if !condition(item, index, slice) {
			return false
		}
	}
	return true
}
