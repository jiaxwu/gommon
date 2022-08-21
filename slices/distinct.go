package slices

// Distinct remove the same items
func Distinct[T comparable](slice []T) []T {
	m := make(map[T]struct{}, len(slice))
	var res []T
	for _, item := range slice {
		if _, ok := m[item]; !ok {
			res = append(res, item)
			m[item] = struct{}{}
		}
	}
	return res
}
