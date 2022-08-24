package math

import (
	"math"

	"golang.org/x/exp/constraints"
)

// 把数n分成m份
// 返回每一份大小和最后一份的大小
func Split(n, m uint) (uint, uint) {
	per := (n + m - 1) / m
	last := n % per
	if last == 0 {
		last = per
	}
	return per, last
}

// 最大值
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// 最小值
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// 绝对值
func Abs[T constraints.Float | constraints.Signed](a T) T {
	if a < 0 {
		return -a
	}
	return a
}

// Log2
func Log2[T constraints.Integer | constraints.Float](a T) T {
	return T(math.Log2(float64(a)))
}
