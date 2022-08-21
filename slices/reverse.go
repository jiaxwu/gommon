package slices

// Reverse slice
func Reverse[T any](slice []T) []T {
	reversed := make([]T, 0, len(slice))
	for i := len(slice) - 1; i >= 0; i-- {
		reversed = append(reversed, slice[i])
	}
	return reversed
}
