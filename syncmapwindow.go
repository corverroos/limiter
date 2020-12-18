package ratelimit

import (
	"sync"
	"sync/atomic"
	"time"
)

// NewSyncMapWindow returns a slightly optimised sync.Map fixed window rate limiter.
// Note: It is not 100% correct, concurrent requests may clobber each other.
// Note: This is not either of the recommended use cases for sync.Map.
func NewSyncMapWindow(period time.Duration, limit int) *SyncMap {
	return &SyncMap{
		period:  period,
		limit:   limit,
		nowFunc: time.Now,
		current: time.Now().Truncate(period).UnixNano(),
		counts:  new(sync.Map),
	}
}

type SyncMap struct {
	period time.Duration
	limit  int

	current int64 // time.Now().Truncate(period).UnixNano() with atomic access only
	counts  *sync.Map

	nowFunc func() time.Time
}

func (n *SyncMap) Request(resource string) bool {
	bucket := n.nowFunc().Truncate(n.period).UnixNano()
	current := atomic.LoadInt64(&n.current)

	if bucket == current {
		var c int
		if v, ok := n.counts.Load(resource); ok {
			c = v.(int)
		}
		c++
		n.counts.Store(resource, c) // Fixed window doesn't distinguish between requests that contribute or not.
		return c <= n.limit
	}

	atomic.StoreInt64(&n.current, bucket)
	n.counts = new(sync.Map)
	n.counts.Store(resource, 1)
	return true
}

var _ RateLimiter = (*SyncMap)(nil)
