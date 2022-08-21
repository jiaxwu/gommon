package slices

// ReduceRight slice to a value
func ReduceRight[T any, R any](slice []T, reduce func(total R, item T, index int, slice []T) R, init R) R {
	for index := len(slice) - 1; index >= 0; index-- {
		init = reduce(init, slice[index], index, slice)
	}
	return init
}
