package slices

// ForEach item execute action
func ForEach[T any](slice []T, action func(item T, index int, slice []T)) {
	for index, item := range slice {
		action(item, index, slice)
	}
}
