package timingwheel

import (
	"context"
	"sync"
	"time"

	"github.com/jiaxwu/gommon/container/heap"
)

// 延迟队列
type delayQueue struct {
	h      *heap.Heap[*bucket]
	mutex  sync.Mutex    // 保证并发安全
	wakeup chan struct{} // 唤醒通道
}

// 创建延迟队列
func newDelayQueue() *delayQueue {
	return &delayQueue{
		h: heap.New(nil, func(b1, b2 *bucket) bool {
			return b1.expiration < b2.expiration
		}),
		wakeup: make(chan struct{}, 1),
	}
}

// 添加延迟元素到队列
func (q *delayQueue) push(b *bucket) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.h.Push(b)
	// 唤醒等待的协程
	// 这里表示新添加的元素到期时间是最早的，或者原来队列为空
	// 因此必须唤醒等待的协程，因为可以拿到更早到期的元素
	if q.h.Peek() == b {
		select {
		case q.wakeup <- struct{}{}:
		default:
		}
	}
}

// 等待直到有元素到期
// 或者ctx被关闭
func (q *delayQueue) take(ctx context.Context) *bucket {
	for {
		var expiration *time.Timer
		q.mutex.Lock()
		// 有元素
		if !q.h.Empty() {
			// 获取元素
			now := time.Now().UnixMilli()
			entry := q.h.Peek()
			if now > entry.expiration {
				q.h.Pop()
				q.mutex.Unlock()
				return entry
			}
			// 到期时间，使用time.NewTimer()才能够调用Stop()，从而释放定时器
			expiration = time.NewTimer(time.Duration(now-entry.expiration) * time.Millisecond)
		}
		// 避免被之前的元素假唤醒
		select {
		case <-q.wakeup:
		default:
		}
		q.mutex.Unlock()

		// 不为空，需要同时等待元素到期，并且除非expiration到期，否则都需要关闭expiration避免泄露
		if expiration != nil {
			select {
			case <-q.wakeup: // 新的更快到期元素
				expiration.Stop()
			case <-expiration.C: // 首元素到期
			case <-ctx.Done(): // 被关闭
				expiration.Stop()
				return nil
			}
		} else {
			select {
			case <-q.wakeup: // 新的更快到期元素
			case <-ctx.Done(): // 被关闭
				return nil
			}
		}
	}
}

// 返回一个通道，输出到期元素
// size是通道缓存大小
func (q *delayQueue) channel(ctx context.Context, size int) <-chan *bucket {
	out := make(chan *bucket, size)
	go func() {
		for {
			entry := q.take(ctx)
			if entry == nil {
				close(out)
				return
			}
			out <- entry
		}
	}()
	return out
}
