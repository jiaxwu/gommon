package slices

// ToMap iterates the slice to map following the keyFunc and valueFunc
func ToMap[T any, K comparable, V any](slice []T, keyFunc func(item T, index int, slice []T) K,
	valueFunc func(item T, index int, slice []T) V) map[K]V {
	dict := make(map[K]V, len(slice))
	for index, value := range slice {
		dict[keyFunc(value, index, slice)] = valueFunc(value, index, slice)
	}
	return dict
}
