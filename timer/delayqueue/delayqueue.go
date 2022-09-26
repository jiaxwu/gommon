package delayqueue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jiaxwu/gommon/container/heap"
)

type entry[T any] struct {
	value      T
	expiration time.Time // 到期时间
}

// 延迟队列
// 参考https://github.com/RussellLuo/timingwheel/blob/master/delayqueue/delayqueue.go
type DelayQueue[T any] struct {
	h *heap.Heap[*entry[T]]
	// // 保证并发安全
	mutex sync.Mutex
	// 表示Take()是否正在等待队列不为空或更早到期的元素
	// 0表示Take()没在等待，1表示Take()在等待
	sleeping int32
	// 唤醒通道
	wakeup chan struct{}
}

// 创建延迟队列
func New[T any]() *DelayQueue[T] {
	return &DelayQueue[T]{
		h: heap.New(nil, func(e1, e2 *entry[T]) bool {
			return e1.expiration.Before(e2.expiration)
		}),
		wakeup: make(chan struct{}),
	}
}

// 添加延迟元素到队列
func (q *DelayQueue[T]) Push(value T, delay time.Duration) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	entry := &entry[T]{
		value:      value,
		expiration: time.Now().Add(delay),
	}
	q.h.Push(entry)
	// 唤醒等待的Take()
	// 这里表示新添加的元素到期时间是最早的，或者原来队列为空
	// 因此必须唤醒等待的Take()，因为可以拿到更早到期的元素
	if q.h.Peek() == entry {
		// 把sleeping从1修改成0，也就是唤醒等待的Take()
		if atomic.CompareAndSwapInt32(&q.sleeping, 1, 0) {
			q.wakeup <- struct{}{}
		}
	}
}

// 等待直到有元素到期
// 或者ctx被关闭
func (q *DelayQueue[T]) Take(ctx context.Context) (T, bool) {
	for {
		var timer *time.Timer
		q.mutex.Lock()
		// 有元素
		if !q.h.Empty() {
			// 获取元素
			entry := q.h.Peek()
			now := time.Now()
			if now.After(entry.expiration) {
				q.h.Pop()
				q.mutex.Unlock()
				return entry.value, true
			}
			// 到期时间，使用time.NewTimer()才能够调用Stop()，从而释放定时器
			timer = time.NewTimer(entry.expiration.Sub(now))
		}
		// 走到这里表示需要等待了，设置为1告诉Push()在有新元素时要通知
		atomic.StoreInt32(&q.sleeping, 1)
		q.mutex.Unlock()

		// 不为空，需要同时等待元素到期，并且除非timer到期，否则都需要关闭timer避免泄露
		if timer != nil {
			select {
			case <-q.wakeup: // 新的更快到期元素
				timer.Stop()
			case <-timer.C: // 首元素到期
				// 设置为0，如果原来也为0表示有Push()正在q.wakeup被阻塞
				if atomic.SwapInt32(&q.sleeping, 0) == 0 {
					// 避免Push()的协程被阻塞
					<-q.wakeup
				}
			case <-ctx.Done(): // 被关闭
				timer.Stop()
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
	if time.Now().Before(entry.expiration) {
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
