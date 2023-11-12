package slices

// First not zero value item
func FirstNotZero[T comparable](items ...T) (T, bool) {
	var t T
	for _, item := range items {
		if item != t {
			return item, true
		}
	}
	return t, false
}
