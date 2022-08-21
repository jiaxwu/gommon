package priority_queue

import "github.com/jiaxwu/gommon/container/heap"

// 优先队列
type PriorityQueue[T any] struct {
	h *heap.Heap[T]
}

func New[T any](h []T, less func(e1 T, e2 T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		h: heap.New(h, less),
	}
}

// 入队
func (p *PriorityQueue[T]) Add(elem T) {
	p.h.Push(elem)
}

// 出队
func (p *PriorityQueue[T]) Remove() T {
	return p.h.Pop()
}

// 队列元素个数
func (p *PriorityQueue[T]) Len() int {
	return p.h.Len()
}

// 队列是否为空
func (p *PriorityQueue[T]) Empty() bool {
	return p.Len() == 0
}
