package mem

import (
	"testing"
)

func TestMemset(t *testing.T) {
	n := 10
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i
	}

	setVals := []int{0, 2, 5, 10}

	for _, setVal := range setVals {
		Memset(arr, setVal)
		for i := 0; i < n; i++ {
			if arr[i] != setVal {
				t.Errorf("want %v, but %v", setVal, arr[i])
			}
		}
	}
}

// 循环设置
func Loopset[T any](arr []T, val T) {
	for i := 0; i < len(arr); i++ {
		arr[i] = val
	}
}

var a = make([]int, 1000)

func BenchmarkLoopset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Loopset(a, 10)
	}
}

func BenchmarkMemset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Memset(a, 10)
	}
}
