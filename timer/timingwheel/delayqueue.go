package timingwheel

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jiaxwu/gommon/container/heap"
)

// 延迟队列
type delayQueue struct {
	h        *heap.Heap[*bucket]
	mutex    sync.Mutex    // 保证并发安全
	sleeping int32         // 用于Push()和Take()之间通知是否有需要唤醒
	wakeup   chan struct{} // 唤醒通道
}

// 创建延迟队列
func newDelayQueue() *delayQueue {
	return &delayQueue{
		h: heap.New(nil, func(b1, b2 *bucket) bool {
			return b1.getExpiration() < b2.getExpiration()
		}),
		wakeup: make(chan struct{}),
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
		if atomic.CompareAndSwapInt32(&q.sleeping, 1, 0) {
			q.wakeup <- struct{}{}
		}
	}
}

// 等待直到有元素到期
// 或者ctx被关闭
func (q *delayQueue) take(ctx context.Context, nowF func() int64) *bucket {
	for {
		var t *time.Timer
		q.mutex.Lock()
		// 有元素
		if !q.h.Empty() {
			// 获取元素
			entry := q.h.Peek()
			expiration := entry.getExpiration()
			now := nowF()
			if now > expiration {
				q.h.Pop()
				q.mutex.Unlock()
				return entry
			}
			// 到期时间，使用time.NewTimer()才能够调用Stop()，从而释放定时器
			t = time.NewTimer(time.Duration(now-expiration) * time.Millisecond)
		}
		// 走到这里表示需要等待了，则需要告诉Push()在有新元素时要通知
		atomic.StoreInt32(&q.sleeping, 1)
		q.mutex.Unlock()

		// 不为空，需要同时等待元素到期，并且除非t到期，否则都需要关闭t避免泄露
		if t != nil {
			select {
			case <-q.wakeup: // 新的更快到期元素
				t.Stop()
			case <-t.C: // 首元素到期
				if atomic.SwapInt32(&q.sleeping, 0) == 0 {
					// 避免Push()的协程被阻塞
					<-q.wakeup
				}
			case <-ctx.Done(): // 被关闭
				t.Stop()
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
func (q *delayQueue) channel(ctx context.Context, size int, nowF func() int64) <-chan *bucket {
	out := make(chan *bucket, size)
	go func() {
		for {
			entry := q.take(ctx, nowF)
			if entry == nil {
				close(out)
				return
			}
			out <- entry
		}
	}()
	return out
}
