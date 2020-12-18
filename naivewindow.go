package ratelimit

import (
	"sync"
	"time"
)

// NewNaiveWindow returns a naive fixed window rate limiter.
// It is naive since uses a single mutex; so no performance optimisations.
func NewNaiveWindow(period time.Duration, limit int) *NaiveWindow {
	return &NaiveWindow{
		period:  period,
		limit:   limit,
		nowFunc: time.Now,
	}
}

type NaiveWindow struct {
	period time.Duration
	limit  int

	current time.Time
	counts  map[string]int
	mu      sync.Mutex

	nowFunc func() time.Time
}

func (n *NaiveWindow) Request(resource string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	bucket := n.nowFunc().Truncate(n.period)
	if bucket == n.current {
		n.counts[resource]++ // Fixed window doesn't distinguish between requests that contribute or not.
		return n.counts[resource] <= n.limit
	}

	n.current = bucket
	n.counts = map[string]int{resource: 1}

	return true
}

var _ RateLimiter = (*NaiveWindow)(nil)
