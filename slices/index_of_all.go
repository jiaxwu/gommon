package slices

// IndexOfAll slice[i] == item
func IndexOfAll[T comparable](slice []T, item T) []int {
	var indexes []int
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			indexes = append(indexes, i)
		}
	}
	return indexes
}
