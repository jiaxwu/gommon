package slices

// Fill items to slice
func Fill[T any](slice []T, itemFunc func(item T, index int, slice []T) T) []T {
	for index := 0; index < len(slice); index++ {
		slice[index] = itemFunc(slice[index], index, slice)
	}
	return slice
}
