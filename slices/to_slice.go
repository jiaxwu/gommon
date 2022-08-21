package slices

// ToSlice iterates map's key and value to generate a new slice
func ToSlice[K comparable, V any, T any](dict map[K]V, transform func(key K, value V) T) []T {
	slice := make([]T, 0, len(dict))
	for key, value := range dict {
		slice = append(slice, transform(key, value))
	}
	return slice
}
