package deque

import "github.com/jiaxwu/gommon/container/list"

type Deque[T any] struct {
	l *list.List[T]
}

func New[T any]() *Deque[T] {
	return &Deque[T]{l: list.New[T]()}
}

// 从队头入队
func (d *Deque[T]) PushFront(elem T) {
	d.l.PushFront(elem)
}

// 从队尾入队
func (d *Deque[T]) PushBack(elem T) {
	d.l.PushBack(elem)
}

// 从队头出队
func (d *Deque[T]) RemoveFront() T {
	return d.l.RemoveFront()
}

// 从队尾出队
func (d *Deque[T]) RemoveBack() T {
	return d.l.RemoveBack()
}

// Len 队列元素个数
func (d *Deque[T]) Len() int {
	return d.l.Len()
}

// Empty 队列是否为空
func (d *Deque[T]) Empty() bool {
	return d.l.Empty()
}
