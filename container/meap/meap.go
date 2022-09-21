package meap

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

type LessFunc[K comparable, V any] func(e1 Entry[K, V], e2 Entry[K, V]) bool

// meap=m[ap]+[h]eap
// 可以log(n)删除任意元素的堆
// 可以支持随机查询
// 是堆和map的结合
// 也就是带有map的特性和堆的特性
type Meap[K comparable, V any] struct {
	h        []Entry[K, V]
	m        map[K]int
	lessFunc LessFunc[K, V]
}

func New[K comparable, V any](lessFunc LessFunc[K, V]) *Meap[K, V] {
	return &Meap[K, V]{
		m:        make(map[K]int),
		lessFunc: lessFunc,
	}
}

// 移除堆顶元素
func (h *Meap[K, V]) Pop() Entry[K, V] {
	n := h.Len() - 1
	h.swap(0, n)
	h.down(0, n)
	return h.pop()
}

// 获取堆顶元素
func (h *Meap[K, V]) Peek() Entry[K, V] {
	return h.h[0]
}

// 获取元素
func (h *Meap[K, V]) Get(key K) (V, bool) {
	index, ok := h.m[key]
	if !ok {
		var v V
		return v, ok
	}
	return h.h[index].Value, true
}

// 添加元素到堆
func (h *Meap[K, V]) Push(key K, value V) {
	// 如果堆中已经包含这个元素
	// 更新值并调整堆
	if h.Contains(key) {
		index := h.m[key]
		h.h[index].Value = value
		h.fix(index)
		return
	}

	// 否则添加元素
	h.push(key, value)
	h.up(h.Len() - 1)
}

// 堆长度
func (h *Meap[K, V]) Len() int {
	return len(h.h)
}

// 堆是否为空
func (h *Meap[K, V]) Empty() bool {
	return h.Len() == 0
}

// 移除堆里对应Key的元素
func (h *Meap[K, V]) Remove(key K) {
	i, ok := h.m[key]
	if !ok {
		return
	}

	n := h.Len() - 1
	if n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	h.pop()
}

// 是否包含这个元素
func (h *Meap[K, V]) Contains(key K) bool {
	_, ok := h.m[key]
	return ok
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *Meap[K, V]) fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}

func (h *Meap[K, V]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *Meap[K, V]) down(i0, n int) bool {
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

func (h *Meap[K, V]) less(i, j int) bool {
	return h.lessFunc(h.h[i], h.h[j])
}

// swap两个元素的时候
// 两个元素在map里的下标也要交换
func (h *Meap[K, V]) swap(i, j int) {
	h.h[i], h.h[j] = h.h[j], h.h[i]
	h.m[h.h[i].Key] = i
	h.m[h.h[j].Key] = j
}

// 添加一个元素到堆的末尾
func (h *Meap[K, V]) push(key K, value V) {
	h.m[key] = h.Len()
	h.h = append(h.h, Entry[K, V]{
		Key:   key,
		Value: value,
	})
}

// 从堆的末尾移除元素
func (h *Meap[K, V]) pop() Entry[K, V] {
	elem := h.h[h.Len()-1]
	h.h = h.h[:h.Len()-1]
	delete(h.m, elem.Key)
	return elem
}
