package slices

// Frequencies returns a map with the unique values of the collection as keys and their frequencies as the values.
func Frequencies[T comparable](slice []T) map[T]int {
	frequencies := map[T]int{}
	for _, item := range slice {
		frequencies[item]++
	}
	return frequencies
}
