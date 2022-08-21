package limiter

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ViolationStrategyError 违背策略错误
type ViolationStrategyError struct {
	Limit  int           // 窗口请求上限
	Window time.Duration // 窗口时间大小
}

func (e *ViolationStrategyError) Error() string {
	return fmt.Sprintf("violation strategy that limit = %d and window = %d", e.Limit, e.Window)
}

// SlidingLogLimiterStrategy 滑动日志限流器的策略
type SlidingLogLimiterStrategy struct {
	limit        int   // 窗口请求上限
	window       int64 // 窗口时间大小
	smallWindows int64 // 小窗口数量
}

func NewSlidingLogLimiterStrategy(limit int, window time.Duration) *SlidingLogLimiterStrategy {
	return &SlidingLogLimiterStrategy{
		limit:  limit,
		window: int64(window),
	}
}

// SlidingLogLimiter 滑动日志限流器
type SlidingLogLimiter struct {
	strategies  []*SlidingLogLimiterStrategy // 滑动日志限流器策略列表
	smallWindow int64                        // 小窗口时间大小
	counters    map[int64]int                // 小窗口计数器
	mutex       sync.Mutex                   // 避免并发问题
}

func NewSlidingLogLimiter(smallWindow time.Duration, strategies ...*SlidingLogLimiterStrategy) (*SlidingLogLimiter, error) {
	// 复制策略避免被修改
	strategies = append(make([]*SlidingLogLimiterStrategy, 0, len(strategies)), strategies...)

	// 不能不设置策略
	if len(strategies) == 0 {
		return nil, errors.New("must be set strategies")
	}

	// 排序策略，窗口时间大的排前面，相同窗口上限大的排前面
	sort.Slice(strategies, func(i, j int) bool {
		a, b := strategies[i], strategies[j]
		if a.window == b.window {
			return a.limit > b.limit
		}
		return a.window > b.window
	})

	for i, strategy := range strategies {
		// 随着窗口时间变小，窗口上限也应该变小
		if i > 0 {
			if strategy.limit >= strategies[i-1].limit {
				return nil, errors.New("the smaller window should be the smaller limit")
			}
		}
		// 窗口时间必须能够被小窗口时间整除
		if strategy.window%int64(smallWindow) != 0 {
			return nil, errors.New("window cannot be split by integers")
		}
		strategy.smallWindows = strategy.window / int64(smallWindow)
	}

	return &SlidingLogLimiter{
		strategies:  strategies,
		smallWindow: int64(smallWindow),
		counters:    make(map[int64]int),
	}, nil
}

func (l *SlidingLogLimiter) TryAcquire() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 获取当前小窗口值
	currentSmallWindow := time.Now().UnixNano() / l.smallWindow * l.smallWindow
	// 获取每个策略的起始小窗口值
	startSmallWindows := make([]int64, len(l.strategies))
	for i, strategy := range l.strategies {
		startSmallWindows[i] = currentSmallWindow - l.smallWindow*(strategy.smallWindows-1)
	}

	// 计算每个策略当前窗口的请求总数
	counts := make([]int, len(l.strategies))
	for smallWindow, counter := range l.counters {
		if smallWindow < startSmallWindows[0] {
			delete(l.counters, smallWindow)
			continue
		}
		for i := range l.strategies {
			if smallWindow >= startSmallWindows[i] {
				counts[i] += counter
			}
		}
	}

	// 若到达对应策略窗口请求上限，请求失败，返回违背的策略
	for i, strategy := range l.strategies {
		if counts[i] >= strategy.limit {
			return &ViolationStrategyError{
				Limit:  strategy.limit,
				Window: time.Duration(strategy.window),
			}
		}
	}

	// 若没到窗口请求上限，当前小窗口计数器+1，请求成功
	l.counters[currentSmallWindow]++
	return nil
}
