package qps

import (
	"log"
	"sync"
	"time"
)

// 窗口信息
type Window struct {
	TotalCnt  int64         // 总次数
	TotalTime time.Duration // 总时间
}

// 平均耗时
func (w *Window) AvgTime() time.Duration {
	return w.TotalTime / time.Duration(w.TotalCnt)
}

// 基于滑动窗口的QPS统计
type QPS struct {
	windowCnt int64            // 窗口数量
	windows   map[int64]Window // 窗口
	mut       sync.Mutex       // 避免并发问题
}

// windowCnt: 一秒分为多少个窗口，越细越准确，但是消耗越大，且必须能够把窗口整除
func New(windowCnt int64) *QPS {
	// 窗口时间必须能够被窗口数量整除
	if time.Second%time.Duration(windowCnt) != 0 {
		log.Fatal("window cannot be split by integers")
	}
	return &QPS{
		windowCnt: windowCnt,
		windows:   make(map[int64]Window),
	}
}

// 记录QPS
func (q *QPS) Add() {
	q.mut.Lock()
	defer q.mut.Unlock()
	q.add(0)
}

// 记录QPS和使用时间
func (q *QPS) AddSince(start time.Time) {
	q.mut.Lock()
	defer q.mut.Unlock()
	q.add(time.Since(start))
}

// 记录QPS和使用时间
func (q *QPS) AddUseTime(useTime time.Duration) {
	q.mut.Lock()
	defer q.mut.Unlock()
	q.add(useTime)
}

func (q *QPS) add(useTime time.Duration) {
	// 当前窗口计数器+1
	curWindowTime := q.curWindowTime()
	window := q.windows[curWindowTime]
	window.TotalCnt++
	window.TotalTime += useTime
	q.windows[curWindowTime] = window
	startWindowTime := q.startWindowTime()
	// 删除过期窗口
	for windowTime := range q.windows {
		if windowTime < startWindowTime {
			delete(q.windows, windowTime)
		}
	}
}

// 获取QPS信息
func (q *QPS) Get() Window {
	q.mut.Lock()
	defer q.mut.Unlock()
	startWindowTime := q.startWindowTime()
	// 计算当前窗口的请求总数
	var w Window
	for windowTime, window := range q.windows {
		if windowTime < startWindowTime {
			delete(q.windows, windowTime)
		} else {
			w.TotalCnt += window.TotalCnt
			w.TotalTime += window.TotalTime
		}
	}
	return w
}

// 窗口时间大小
func (q *QPS) WindowSize() time.Duration {
	return time.Second / time.Duration(q.windowCnt)
}

// 当前窗口时间
func (q *QPS) curWindowTime() int64 {
	windowSize := int64(q.WindowSize())
	return time.Now().UnixNano() / windowSize * windowSize
}

// 起始窗口时间
func (q *QPS) startWindowTime() int64 {
	windowSize := int64(q.WindowSize())
	return time.Now().UnixNano()/windowSize*windowSize - windowSize*(q.windowCnt-1)
}
