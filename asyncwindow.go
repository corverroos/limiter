package ratelimit

import (
	"sync/atomic"
	"time"
)

const asyncChs = 1

// AsyncWindow returns a optimised async fixed window rate limiter.
// Note: It is not 100% correct, requests may be dropped and it lags by one window.
func NewAsyncWindow(period time.Duration, limit int) *AsyncWindow {
	a := &AsyncWindow{
		period: period,
		limit:  limit,
		ch:     make(chan string, 10000),
	}
	go a.runForever()

	return a
}

type AsyncWindow struct {
	limit  int
	period time.Duration

	prev    map[string]int
	next    map[string]int
	dropped int64
	current time.Time
	ch      chan string
}

func (a *AsyncWindow) runForever() {
	for resource := range a.ch {
		bucket := time.Now().Truncate(a.period)
		if a.current != bucket {
			a.prev = a.next
			a.next = make(map[string]int)
			a.current = bucket
		}
		a.next[resource]++
	}
}

func (a *AsyncWindow) Request(resource string) bool {
	select {
	case a.ch <- resource:
	default:
		atomic.AddInt64(&a.dropped, 1)
	}

	return a.prev[resource] < a.limit
}
