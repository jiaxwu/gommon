package math

import (
	"math"
	"math/bits"

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
// 修改更加高效算法
func Log2[T constraints.Integer | constraints.Float](a T) T {
	return T(math.Log2(float64(a)))
}

// IsPowOf2 check if a value is a power of two
// Determine whether some value is a power of two, where zero is
// not considered a power of two.
func IsPowOf2[T constraints.Unsigned](n T) bool {
	return n != 0 && ((n & (n - 1)) == 0)
}

// RoundUpPowOf2 round up to nearest power of two
func RoundUpPowOf2[T constraints.Unsigned](n T) T {
	return 1 << FindLastBitSet(n-1)
}

// RoundDownPowOf2 round down to nearest power of two
func RoundDownPowOf2[T constraints.Unsigned](n T) T {
	return 1 << (FindLastBitSet(n) - 1)
}

// FindLastBitSet find last (most-significant) bit set
// This is defined the same way as ffs.
// Note FindLastBitSet(0) = 0, FindLastBitSet(1) = 1, FindLastBitSet(0x80000000) = 32.
func FindLastBitSet[T constraints.Unsigned](x T) int {
	return bits.Len64(uint64(x))
}
