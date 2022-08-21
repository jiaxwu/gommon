package slices

// Equal slice1 equal slice2
func Equal[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for index, item := range slice1 {
		if item != slice2[index] {
			return false
		}
	}
	return true
}

// EqualFunc slice1 equal slice2 with func
func EqualFunc[T any](slice1, slice2 []T, f func(item1, item2 T) bool) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for index, item := range slice1 {
		if !f(item, slice2[index]) {
			return false
		}
	}
	return true
}
