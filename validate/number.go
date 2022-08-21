package validate

import "golang.org/x/exp/constraints"

// 数值在范围内
func Range[T constraints.Ordered](n, min, max T, errFunc func(min, max T) error) error {
	if n < min || n > max {
		return errFunc(min, max)
	}
	return nil
}
