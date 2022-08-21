package limiter

import (
	"sync"
	"time"
)

// LeakyBucketLimiter 漏桶限流器
type LeakyBucketLimiter struct {
	peakLevel       int        // 最高水位
	currentLevel    int        // 当前水位
	currentVelocity int        // 水流速度/秒
	lastTime        time.Time  // 上次放水时间
	mutex           sync.Mutex // 避免并发问题
}

func NewLeakyBucketLimiter(peakLevel, currentVelocity int) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		peakLevel:       peakLevel,
		currentVelocity: currentVelocity,
		lastTime:        time.Now(),
	}
}

func (l *LeakyBucketLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 尝试放水
	now := time.Now()
	// 距离上次放水的时间
	interval := now.Sub(l.lastTime)
	if interval >= time.Second {
		// 当前水位-距离上次放水的时间(秒)*水流速度
		l.currentLevel = maxInt(0, l.currentLevel-int(interval/time.Second)*l.currentVelocity)
		l.lastTime = now
	}

	// 若到达最高水位，请求失败
	if l.currentLevel >= l.peakLevel {
		return false
	}
	// 若没有到达最高水位，当前水位+1，请求成功
	l.currentLevel++
	return true
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
