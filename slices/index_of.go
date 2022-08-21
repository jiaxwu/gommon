package slices

// IndexOf slice[i] == item
func IndexOf[T comparable](slice []T, item T) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			return i
		}
	}
	return -1
}
