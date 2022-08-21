package slices

// FindLast item that meet the condition
func FindLast[T any](slice []T, condition func(item T, index int, slice []T) bool) (T, bool) {
	for i := len(slice) - 1; i >= 0; i-- {
		if condition(slice[i], i, slice) {
			return slice[i], true
		}
	}
	var t T
	return t, false
}
