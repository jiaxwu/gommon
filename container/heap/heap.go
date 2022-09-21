package heap

type LessFunc[T any] func(e1 T, e2 T) bool

type Heap[T any] struct {
	h        []T
	lessFunc LessFunc[T]
}

func New[T any](h []T, lessFunc LessFunc[T]) *Heap[T] {
	heap := &Heap[T]{
		h:        h,
		lessFunc: lessFunc,
	}
	heap.init()
	return heap
}

// 移除堆顶元素
func (h *Heap[T]) Pop() T {
	n := h.Len() - 1
	h.swap(0, n)
	h.down(0, n)
	return h.pop()
}

// 获取堆顶元素
func (h *Heap[T]) Peek() T {
	return h.h[0]
}

// 添加元素到堆
func (h *Heap[T]) Push(x T) {
	h.push(x)
	h.up(h.Len() - 1)
}

// 堆长度
func (h *Heap[T]) Len() int {
	return len(h.h)
}

// 堆是否为空
func (h *Heap[T]) Empty() bool {
	return h.Len() == 0
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[T]) Remove(i int) T {
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
func (h *Heap[T]) Fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Len().
func (h *Heap[T]) init() {
	// heapify
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

func (h *Heap[T]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *Heap[T]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

func (h *Heap[T]) less(i, j int) bool {
	return h.lessFunc(h.h[i], h.h[j])
}

func (h *Heap[T]) swap(i, j int) {
	h.h[i], h.h[j] = h.h[j], h.h[i]
}

func (h *Heap[T]) push(x T) {
	h.h = append(h.h, x)
}

func (h *Heap[T]) pop() T {
	elem := h.h[h.Len()-1]
	h.h = h.h[:h.Len()-1]
	return elem
}
