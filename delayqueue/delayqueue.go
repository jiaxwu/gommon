package delayqueue

import (
	"context"
	"sync"
	"time"

	"github.com/jiaxwu/gommon/container/meap"
)

type Entry[K comparable, V any] struct {
	Key     K
	Value   V
	Expired time.Time // 过期时间
}

// 延迟队列
type DelayQueue[K comparable, V any] struct {
	meap   *meap.Meap[K, Entry[K, V]]
	mutex  sync.Mutex
	wakeup chan struct{} // 唤醒通道
}

func New[K comparable, V any]() *DelayQueue[K, V] {
	return &DelayQueue[K, V]{
		meap: meap.New(func(e1, e2 meap.Entry[K, Entry[K, V]]) bool {
			return e1.Value.Expired.Before(e2.Value.Expired)
		}),
		wakeup: make(chan struct{}, 1),
	}
}

// 添加延迟元素到队列
func (q *DelayQueue[K, V]) Push(key K, value V, delay time.Duration) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.meap.Push(key, Entry[K, V]{
		Key:     key,
		Value:   value,
		Expired: time.Now().Add(delay),
	})
	// 唤醒等待的协程
	if q.meap.Peek().Key == key {
		select {
		case q.wakeup <- struct{}{}:
		default:
		}
	}
}

// 等待直到有元素过期
// 或者ctx被关闭
func (q *DelayQueue[K, V]) Take(ctx context.Context) (Entry[K, V], bool) {
	for {
		q.mutex.Lock()
		var expired <-chan time.Time
		// 有元素
		if !q.meap.Empty() {
			// 获取元素
			entry := q.meap.Peek()
			if time.Now().After(entry.Value.Expired) {
				q.meap.Pop()
				q.mutex.Unlock()
				return entry.Value, true
			}
			// 过期时间
			expired = time.After(time.Until(entry.Value.Expired))
		}
		// 避免被之前的元素假唤醒
		select {
		case <-q.wakeup:
		default:
		}
		q.mutex.Unlock()

		select {
		// 新的更快过期元素
		case <-q.wakeup:
			// 首元素过期
		case <-expired:
			// 被关闭
		case <-ctx.Done():
			return Entry[K, V]{}, false
		}
	}
}

// 返回一个通道，输出过期元素
// size是通道缓存大小
func (q *DelayQueue[K, V]) Channel(ctx context.Context, size int) <-chan Entry[K, V] {
	out := make(chan Entry[K, V], size)
	go func() {
		for {
			entry, ok := q.Take(ctx)
			if !ok {
				return
			}
			out <- entry
		}
	}()
	return out
}

// 获取队头元素
func (q *DelayQueue[K, V]) Peek() (Entry[K, V], bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.meap.Empty() {
		return Entry[K, V]{}, false
	}
	return q.meap.Peek().Value, true
}

// 获取过期元素
func (q *DelayQueue[K, V]) Pop() (Entry[K, V], bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	// 没元素
	if q.meap.Empty() {
		return Entry[K, V]{}, false
	}
	entry := q.meap.Peek()
	// 还没元素过期
	if time.Now().Before(entry.Value.Expired) {
		return Entry[K, V]{}, false
	}
	// 移除元素
	q.meap.Pop()
	return entry.Value, true
}

// 获取元素
func (q *DelayQueue[K, V]) Get(key K) (Entry[K, V], bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.meap.Get(key)
}

// 移除任务
func (q *DelayQueue[K, V]) Remove(key K) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.meap.Remove(key)
}

// 是否队列为空
func (q *DelayQueue[K, V]) Empty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.meap.Empty()
}
