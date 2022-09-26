package timingwheel

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/jiaxwu/gommon/container/list"
)

// 返回给用户的定时器
type Timer struct {
	expiration int64                 // 到期时间
	task       func()                // 任务
	b          unsafe.Pointer        // 所属时间轮的桶
	elem       *list.Element[*Timer] // 为了能从链表中删除
}

func (t *Timer) Stop() bool {
	stoped := false
	for b := t.getBucket(); b != nil; b = t.getBucket() {
		stoped = b.remove(t)
	}
	return stoped
}

func (t *Timer) getBucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.b))
}

func (t *Timer) setBucket(b *bucket) {
	atomic.StorePointer(&t.b, unsafe.Pointer(b))
}

// 时间轮上的一个桶
// 会带着一组相同范围过期时间的元素
type bucket struct {
	expiration int64              // 到期时间
	timers     *list.List[*Timer] // 定时器列表
	mutex      sync.Mutex
}

// 创建定时器
func newBucket() *bucket {
	return &bucket{
		expiration: -1,
		timers:     list.New[*Timer](),
	}
}

// 添加定时器
// 会设置定时器的elem和bucket
// 表示这个定时器所属的桶
func (b *bucket) add(t *Timer) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	elem := b.timers.PushBack(t)
	t.elem = elem
	t.setBucket(b)
}

// 移除定时器
func (b *bucket) remove(t *Timer) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if t.getBucket() != b {
		return false
	}
	b.timers.Remove(t.elem)
	t.setBucket(nil)
	t.elem = nil
	return true
}

// 添加到上一级定时器或执行任务
func (b *bucket) flush(addOrRun func(t *Timer)) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for elem := b.timers.Front(); elem != nil; {
		next := elem.Next()
		t := elem.Value
		if t.getBucket() == b {
			t.setBucket(nil)
			t.elem = nil
		}
		addOrRun(t)
		elem = next
	}

	// 设置过期时间表示没有加入到延迟队列
	b.setExpiration(-1)
	b.timers.Clear()
}

func (b *bucket) getExpiration() int64 {
	return atomic.LoadInt64(&b.expiration)
}

// 返回true表示设置成功
// 否则表示没变化
func (b *bucket) setExpiration(expiration int64) bool {
	return atomic.SwapInt64(&b.expiration, expiration) != expiration
}
