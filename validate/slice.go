package validate

// 长度在范围内
func SliceLen[T any](s []T, min, max int, errFunc func(min, max int) error) error {
	if len(s) < min || len(s) > max {
		return errFunc(min, max)
	}
	return nil
}
