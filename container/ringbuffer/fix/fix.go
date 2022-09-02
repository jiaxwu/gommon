package fix

import (
	"math"

	mmath "github.com/jiaxwu/gommon/math"
)

// 最大长度
const MaxSize = math.MaxInt64

// 固定长度
// 单生产者和单消费者情况下是线程安全性的，但是不能用Reset()方法
type Ring[T any] struct {
	in   uint64 // 写索引
	out  uint64 // 读索引
	size uint64 // 长度
	data []T    // 数据
}

func New[T any](size uint64) *Ring[T] {
	if size == 0 {
		panic("size must be greater than 0")
	}
	if size > MaxSize {
		panic("size is too large")
	}

	return &Ring[T]{
		size: size,
		data: make([]T, size),
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
		panic("ring full")
	}
	in := r.in % r.size
	r.in++
	r.data[in] = e
}

// 写入队尾
func (r *Ring[T]) MPush(elems ...T) {
	size := uint64(len(elems))
	// 不能大于剩余长度
	if size > r.Avail() {
		size = r.Avail()
		elems = elems[:size]
	}
	if size == 0 {
		return
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
