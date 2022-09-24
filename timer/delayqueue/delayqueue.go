package delayqueue

import (
	"context"
	"sync"
	"time"

	"github.com/jiaxwu/gommon/container/heap"
)

type entry[T any] struct {
	value   T
	expired time.Time // 到期时间
}

// 延迟队列
type DelayQueue[T any] struct {
	h      *heap.Heap[*entry[T]]
	mutex  sync.Mutex    // 保证并发安全
	wakeup chan struct{} // 唤醒通道
}

// 创建延迟队列
func New[T any]() *DelayQueue[T] {
	return &DelayQueue[T]{
		h: heap.New(nil, func(e1, e2 *entry[T]) bool {
			return e1.expired.Before(e2.expired)
		}),
		wakeup: make(chan struct{}, 1),
	}
}

// 添加延迟元素到队列
func (q *DelayQueue[T]) Push(value T, delay time.Duration) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	entry := &entry[T]{
		value:   value,
		expired: time.Now().Add(delay),
	}
	q.h.Push(entry)
	// 唤醒等待的协程
	// 这里表示新添加的元素到期时间是最早的，或者原来队列为空
	// 因此必须唤醒等待的协程，因为可以拿到更早到期的元素
	if q.h.Peek() == entry {
		select {
		case q.wakeup <- struct{}{}:
		default:
		}
	}
}

// 等待直到有元素到期
// 或者ctx被关闭
func (q *DelayQueue[T]) Take(ctx context.Context) (T, bool) {
	for {
		var expired *time.Timer
		q.mutex.Lock()
		// 有元素
		if !q.h.Empty() {
			// 获取元素
			entry := q.h.Peek()
			if time.Now().After(entry.expired) {
				q.h.Pop()
				q.mutex.Unlock()
				return entry.value, true
			}
			// 到期时间，使用time.NewTimer()才能够调用Stop()，从而释放定时器
			expired = time.NewTimer(time.Until(entry.expired))
		}
		// 避免被之前的元素假唤醒
		select {
		case <-q.wakeup:
		default:
		}
		q.mutex.Unlock()

		// 不为空，需要同时等待元素到期，并且除非expired到期，否则都需要关闭expired避免泄露
		if expired != nil {
			select {
			case <-q.wakeup: // 新的更快到期元素
				expired.Stop()
			case <-expired.C: // 首元素到期
			case <-ctx.Done(): // 被关闭
				expired.Stop()
				var t T
				return t, false
			}
		} else {
			select {
			case <-q.wakeup: // 新的更快到期元素
			case <-ctx.Done(): // 被关闭
				var t T
				return t, false
			}
		}
	}
}

// 返回一个通道，输出到期元素
// size是通道缓存大小
func (q *DelayQueue[T]) Channel(ctx context.Context, size int) <-chan T {
	out := make(chan T, size)
	go func() {
		for {
			entry, ok := q.Take(ctx)
			if !ok {
				close(out)
				return
			}
			out <- entry
		}
	}()
	return out
}

// 获取队头元素
func (q *DelayQueue[T]) Peek() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.h.Empty() {
		var t T
		return t, false
	}
	return q.h.Peek().value, true
}

// 获取到期元素
func (q *DelayQueue[T]) Pop() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	// 没元素
	if q.h.Empty() {
		var t T
		return t, false
	}
	entry := q.h.Peek()
	// 还没元素到期
	if time.Now().Before(entry.expired) {
		var t T
		return t, false
	}
	// 移除元素
	q.h.Pop()
	return entry.value, true
}

// 是否队列为空
func (q *DelayQueue[T]) Empty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.h.Empty()
}
