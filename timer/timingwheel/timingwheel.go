package timingwheel

import (
	"context"
	"time"
	"unsafe"
)

const (
	tick                 = 1  // 每一跳时间
	wheelSize            = 20 // 跳数
	delayQueueBufferSize = 10 // 延迟队列缓冲区大小
)

// 时间轮
// 单位都是毫秒
type TimingWheel struct {
	tick        int64          // 每一跳的时间
	wheelSize   int64          // 时间轮
	interval    int64          // 一圈的时间
	currentTime int64          // 当前时间
	buckets     []*bucket      // 时间轮的每个桶
	queue       *delayQueue    // 桶延迟队列
	prevWheel   *TimingWheel   // 上一个时间轮
	nextWheel   unsafe.Pointer // 下一个时间轮
}

func New() *TimingWheel {
	return newTimingWheel(tick, newDelayQueue(), time.Now().UnixMilli(), nil)
}

func newTimingWheel(tick int64, queue *delayQueue, currentTime int64 prevWheel *TimingWheel) *TimingWheel {
	timingWheel := &TimingWheel{
		tick:        tick,
		wheelSize:   wheelSize,
		interval:    tick * wheelSize,
		currentTime: currentTime,
		buckets:     make([]*bucket, wheelSize),
		queue:       queue,
		prevWheel:   prevWheel,
	}
	for i := 0; i < wheelSize; i++ {
		timingWheel.buckets[i] = newBucket()
	}
	return timingWheel
}

// 运行时间轮
func (tw *TimingWheel) Run(ctx context.Context) {
	bucketChan := tw.queue.channel(ctx, delayQueueBufferSize)
	for {
		select {
		case b, ok := <-bucketChan: // 桶到期
			if !ok {
				return
			}
			// 前进当前时间
			tw.advance(b.expiration)
			// 处理桶
			tw.flush(b)
		case timer := <-tw.addChan: // 添加元素
			tw.add(timer)
		case timer := <-tw.removeChan: // 删除元素
			tw.remove(timer)
		case <-ctx.Done(): // 被关闭
			return
		}
	}
}

// 添加定时器
func (tw *TimingWheel) AfterFunc(delay time.Duration, f func()) *Timer {
	t := &Timer{
		expiration: time.Now().Add(delay).UnixMilli(),
		f:          f,
	}
	tw.addChan <- t
	return t
}

// 删除定时器
// 不要重复删除，否则会产生panic
func (tw *TimingWheel) Remove(t *Timer) {
	if t == nil {
		return
	}
	tw.removeChan <- t
}

// 处理桶到期任务
func (tw *TimingWheel) flush(b *bucket) {
	if tw.prevWheel == nil { // 第一级时间轮
		for elem := b.timers.Front(); elem != nil; elem = elem.Next() {
			t := elem.Value
			t.bucket = nil
			t.elem = nil
			go t.f()
		}
	} else { // 其他时间轮
		for elem := b.timers.Front(); elem != nil; elem = elem.Next() {
			// 添加到上一级时间轮
			tw.prevWheel.add(elem.Value)
		}
	}
	// 清空桶
	b.timers.Clear()
}

// 添加定时器
func (tw *TimingWheel) add(t *Timer) {
	if t.expiration < tw.currentTime+tw.tick { // 已经过期了
		go t.f()
	} else if t.expiration < tw.currentTime+tw.interval { // 在当前时间轮里
		// 多少跳了
		ticks := t.expiration / tw.tick
		b := tw.buckets[ticks%tw.wheelSize]
		b.add(t)
		// 第一个加入桶，因此bucket还没入队
		if b.timers.Len() == 1 {
			// 设置桶到期时间
			b.expiration = ticks * tw.tick
			tw.queue.push(b)
		}
	} else { // 在其他时间轮里
		if tw.nextWheel == nil {
			tw.nextWheel = newTimingWheel(tw.interval, tw.queue, nil)
		}
		tw.nextWheel.add(t)
	}
}

// 前进时间
func (tw *TimingWheel) advance(expiration int64) {
	if expiration >= tw.currentTime+tw.tick {
		currentTime := expiration - expiration%tw.tick
		tw.currentTime = currentTime
		if tw.nextWheel != nil {
			tw.nextWheel.advance(currentTime)
		}
	}
}

// 删除定时器
func (tw *TimingWheel) remove(t *Timer) {
	if t.bucket == nil || t.elem == nil {
		return
	}
	t.bucket.remove(t)
}
