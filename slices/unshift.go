package slices

// Unshift add items from head
func Unshift[T any](slice []T, items ...T) []T {
	return append(items, slice...)
}
