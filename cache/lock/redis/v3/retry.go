package v3

import (
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

type RetryStrategy interface {
	// Next 返回下一次重试的间隔，如果不需要继续重试，那么第二个参数返回 false
	Next() (time.Duration, bool)
	// 以下是几种变体
	// Next(error) (time.Duration, error)
	// Next(context.Context) (time.Duration, error)
	// Next(context.Context,error) (time.Duration, error)
}

// FixIntervalRetryStrategy 固定时间间隔的重试
type FixIntervalRetryStrategy struct {
	Interval time.Duration
	Max      int
	cnt      int
}

func (s *FixIntervalRetryStrategy) Next() (time.Duration, bool) {
	s.cnt++
	return s.Interval, s.cnt <= s.Max
}

// BackOffRetryStrategy 指数退避
type BackOffRetryStrategy struct {
	factor   float64
	minDelay time.Duration
	maxDelay time.Duration
	attempts uint64
	jitter   bool
}

func NewDefaultBackOffRetryStrategy() *BackOffRetryStrategy {
	return &BackOffRetryStrategy{
		factor:   2,
		minDelay: 100 * time.Millisecond,
		maxDelay: 2 * time.Second,
	}
}

func (b *BackOffRetryStrategy) Next() (time.Duration, bool) {
	dur := float64(b.minDelay) * math.Pow(b.factor, float64(b.attempts))
	if b.jitter == true {
		dur = rand.Float64()*(dur-float64(b.minDelay)) + float64(b.minDelay)
	}
	if dur > float64(b.maxDelay) {
		dur = float64(b.maxDelay)
	}
	atomic.AddUint64(&b.attempts, 1)
	return time.Duration(dur), true
}
