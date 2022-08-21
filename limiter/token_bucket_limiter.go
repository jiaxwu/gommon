package limiter

import (
	"sync"
	"time"
)

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	capacity      int        // 容量
	currentTokens int        // 令牌数量
	rate          int        // 发放令牌速率/秒
	lastTime      time.Time  // 上次发放令牌时间
	mutex         sync.Mutex // 避免并发问题
}

func NewTokenBucketLimiter(capacity, rate int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity: capacity,
		rate:     rate,
		lastTime: time.Now(),
	}
}

func (l *TokenBucketLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 尝试发放令牌
	now := time.Now()
	// 距离上次发放令牌的时间
	interval := now.Sub(l.lastTime)
	if interval >= time.Second {
		// 当前令牌数量+距离上次发放令牌的时间(秒)*发放令牌速率
		l.currentTokens = minInt(l.capacity, l.currentTokens+int(interval/time.Second)*l.rate)
		l.lastTime = now
	}

	// 如果没有令牌，请求失败
	if l.currentTokens == 0 {
		return false
	}
	// 如果有令牌，当前令牌-1，请求成功
	l.currentTokens--
	return true
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
