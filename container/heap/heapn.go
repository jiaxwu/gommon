package heap

// n叉堆
type HeapN[T any] struct {
	n        int
	h        []T
	lessFunc LessFunc[T]
}

func NewN[T any](n int, lessFunc LessFunc[T]) *HeapN[T] {
	if n < 2 {
		panic("n must be greater than 2")
	}
	return &HeapN[T]{
		n:        n,
		lessFunc: lessFunc,
	}
}

// 移除堆顶元素
func (h *HeapN[T]) Pop() T {
	n := h.Len() - 1
	h.swap(0, n)
	h.down(0, n)
	return h.pop()
}

// 获取堆顶元素
func (h *HeapN[T]) Peek() T {
	return h.h[0]
}

// 添加元素到堆
func (h *HeapN[T]) Push(x T) {
	h.push(x)
	h.up(h.Len() - 1)
}

// 堆长度
func (h *HeapN[T]) Len() int {
	return len(h.h)
}

// 堆是否为空
func (h *HeapN[T]) Empty() bool {
	return h.Len() == 0
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *HeapN[T]) Remove(i int) T {
	n := h.Len() - 1
	if n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return h.pop()
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *HeapN[T]) Fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}

func (h *HeapN[T]) up(j int) {
	for {
		i := (j - 1) / h.n // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *HeapN[T]) down(i0, n int) bool {
	i := i0
	for {
		k := h.n*i + 1
		if k >= n || k < 0 { // j1 < 0 after int overflow
			break
		}
		j := k // first child
		for j1 := k + 1; j1 < k+4 && j1 < n; j1++ {
			if h.less(j1, j) {
				j = j1
			}
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

func (h *HeapN[T]) less(i, j int) bool {
	return h.lessFunc(h.h[i], h.h[j])
}

func (h *HeapN[T]) swap(i, j int) {
	h.h[i], h.h[j] = h.h[j], h.h[i]
}

func (h *HeapN[T]) push(x T) {
	h.h = append(h.h, x)
}

func (h *HeapN[T]) pop() T {
	elem := h.h[h.Len()-1]
	h.h = h.h[:h.Len()-1]
	return elem
}
