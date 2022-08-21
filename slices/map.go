package slices

// Map []T1 to []T2 by mapper
func Map[T1 any, T2 any](slice []T1, mapper func(item T1, index int, slice []T1) T2) []T2 {
	mapped := make([]T2, 0, len(slice))
	for index, item := range slice {
		mapped = append(mapped, mapper(item, index, slice))
	}
	return mapped
}
