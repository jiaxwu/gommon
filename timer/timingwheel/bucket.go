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
	f          func()                // 任务
	bucket     unsafe.Pointer        // 所属时间轮的桶
	elem       *list.Element[*Timer] // 为了能从链表中删除
}

func (t *Timer) getBucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.bucket))
}

func (t *Timer) setBucket(b *bucket) {
	atomic.StorePointer(&t.bucket, unsafe.Pointer(b))
}

func (t *Timer) Stop() bool {
	stoped := false
	for b := t.getBucket(); b != nil; b = t.getBucket() {
		stoped = b.remove(t)
	}
	return stoped
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
		timers: list.New[*Timer](),
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
	t.elem = nil
	t.setBucket(nil)
	return true
}

// 循环处理定时器列表
// f可以是执行定时器或者添加到上一级时间轮
func (b *bucket) flush(f func(t *Timer)) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for elem := b.timers.Front(); elem != nil; elem = elem.Next() {
		t := elem.Value
		t.bucket = nil
		t.elem = nil
		f(t)
	}
	b.timers.Clear()
}
