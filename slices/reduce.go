package slices

// Reduce slice to a value
func Reduce[T any, R any](slice []T, reduce func(total R, item T, index int, slice []T) R, init R) R {
	for index := 0; index < len(slice); index++ {
		init = reduce(init, slice[index], index, slice)
	}
	return init
}
