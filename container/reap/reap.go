package meap

type Entry[T any] struct {
	value T
	index int
	h     *Reap[T]
}

func (e *Entry[T]) Value() T {
	return e.value
}

// 移除堆里对应Key的元素
func (e *Entry[T]) Remove() {
	e.h.Remove(e)
}

type LessFunc[T any] func(e1, e2 T) bool

// reap=r[emovable]+[h]eap
// 可以通过Entry实现log(n)删除任意元素的堆
type Reap[T any] struct {
	h        []*Entry[T]
	lessFunc LessFunc[T]
}

func New[T any](lessFunc LessFunc[T]) *Reap[T] {
	return &Reap[T]{
		lessFunc: lessFunc,
	}
}

// 移除堆顶元素
func (h *Reap[T]) Pop() T {
	n := h.Len() - 1
	h.swap(0, n)
	h.down(0, n)
	return h.pop()
}

// 获取堆顶元素
func (h *Reap[T]) Peek() T {
	return h.h[0].value
}

// 添加元素到堆
func (h *Reap[T]) Push(value T) *Entry[T] {
	// 否则添加元素
	entry := h.push(value)
	h.up(h.Len() - 1)
	return entry
}

// 移除堆里对应index的元素
func (h *Reap[T]) Remove(e *Entry[T]) {
	// 不能已经被删除或者不是这个堆所属元素
	if e.index == -1 || e.h != h {
		return
	}
	i := e.index
	n := h.Len() - 1
	if n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	h.pop()
}

// 堆长度
func (h *Reap[T]) Len() int {
	return len(h.h)
}

// 堆是否为空
func (h *Reap[T]) Empty() bool {
	return h.Len() == 0
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *Reap[T]) fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}

func (h *Reap[T]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *Reap[T]) down(i0, n int) bool {
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

func (h *Reap[T]) less(i, j int) bool {
	return h.lessFunc(h.h[i].value, h.h[j].value)
}

// swap两个元素的时候
// 两个元素在map里的下标也要交换
func (h *Reap[T]) swap(i, j int) {
	h.h[i], h.h[j] = h.h[j], h.h[i]
	h.h[i].index = i
	h.h[j].index = j
}

// 添加一个元素到堆的末尾
func (h *Reap[T]) push(value T) *Entry[T] {
	entry := &Entry[T]{
		value: value,
		index: h.Len(),
		h:     h,
	}
	h.h = append(h.h, entry)
	return entry
}

// 从堆的末尾移除元素
func (h *Reap[T]) pop() T {
	elem := h.h[h.Len()-1]
	h.h = h.h[:h.Len()-1]
	elem.h = nil
	elem.index = -1
	return elem.value
}
