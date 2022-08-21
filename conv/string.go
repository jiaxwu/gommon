package conv

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

// 数值转字符串
func Itoa[T constraints.Integer](i T) string {
	return strconv.Itoa(int(i))
}

// 字符串转数值
func Atoi[T constraints.Integer](a string) (T, error) {
	i, err := strconv.Atoi(a)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}
