package unbounded

import (
	"math"

	mmath "github.com/jiaxwu/gommon/math"
)

// 最大长度
const MaxSize = math.MaxInt64

// 动态扩展长度
// 非线程安全，请加锁
type Ring[T any] struct {
	in   uint64 // 写索引
	out  uint64 // 读索引
	size uint64 // 长度
	data []T    // 数据
}

func New[T any](initSize uint64) *Ring[T] {
	if initSize > MaxSize {
		panic("size is too large")
	}

	return &Ring[T]{
		size: initSize,
		data: make([]T, initSize),
	}
}

// 弹出队头元素
func (r *Ring[T]) Pop() T {
	if r.Empty() {
		panic("ring emtpy")
	}
	out := r.out % r.size
	r.out++
	return r.data[out]
}

// 队头元素
func (r *Ring[T]) Peek() T {
	if r.Empty() {
		panic("ring emtpy")
	}
	return r.data[r.out%r.size]
}

// 插入元素到队尾
func (r *Ring[T]) Push(e T) {
	if r.Full() {
		r.Grow(r.Cap() + 1)
	}
	in := r.in % r.size
	r.in++
	r.data[in] = e
}

// 写入队尾
func (r *Ring[T]) MPush(elems ...T) {
	size := uint64(len(elems))
	if size == 0 {
		return
	}
	if size > r.Avail() {
		r.Grow(r.Cap() + size)
	}
	in := r.in % r.size
	copied := copy(r.data[in:], elems)
	copy(r.data, elems[copied:])
	r.in += size
}

// 从队头读取
func (r *Ring[T]) MPop(size uint64) []T {
	if size > r.Len() {
		size = r.Len()
	}
	if size == 0 {
		return nil
	}
	out := r.out % r.size
	elems := make([]T, size)
	copied := copy(elems, r.data[out:])
	copy(elems[copied:], r.data)
	r.out += size
	return elems
}

// 从队头读取，填充到dst里
func (r *Ring[T]) MPopCopy(dst []T) {
	out := r.out % r.size
	dst = dst[:mmath.Min(uint64(len(dst)), r.Len())]
	copied := copy(dst, r.data[out:])
	copied += copy(dst[copied:], r.data)
	r.out += uint64(copied)
}

// 重置读写指针
func (r *Ring[T]) Reset() {
	r.in = 0
	r.out = 0
	r.data = make([]T, r.size)
}

// 总长度
func (r *Ring[T]) Cap() uint64 {
	return r.size
}

// 使用长度
func (r *Ring[T]) Len() uint64 {
	return r.in - r.out
}

// 可用长度
func (r *Ring[T]) Avail() uint64 {
	return r.Cap() - r.Len()
}

// 是否为空
func (r *Ring[T]) Empty() bool {
	return r.in == r.out
}

// 是否满了
func (r *Ring[T]) Full() bool {
	return r.Avail() == 0
}

// 扩容
func (r *Ring[T]) Grow(minSize uint64) {
	size := mmath.Max(r.size*2, minSize)
	if size > MaxSize {
		panic("size is too large")
	}
	if size < 2 {
		size = 2
	}
	// 还没容量，直接申请，因为不需要迁移元素
	if r.size == 0 {
		r.data = make([]T, size)
		r.size = size
		return
	}
	data := make([]T, size)
	out := r.out % r.size
	len := r.Len()
	copied := copy(data[:len], r.data[out:])
	copy(data[copied:len], r.data)
	r.out = 0
	r.in = len
	r.size = size
	r.data = data
}
