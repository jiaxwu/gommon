package slices

// Shift remove first item
func Shift[T any](slice []T) ([]T, T) {
	t := slice[0]
	return slice[1:], t
}
