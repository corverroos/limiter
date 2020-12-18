package ratelimit_test

import (
	"sync"
	"testing"
	"time"

	"github.com/corverroos/ratelimit"
	"github.com/stretchr/testify/require"
)

func TestAsyncWindow(t *testing.T) {
	l := ratelimit.NewAsyncWindow(time.Hour, 10)
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			require.True(t, l.Request(""))
			wg.Done()
		}()
	}
	wg.Wait()
	require.True(t, l.Request(""))
}

func BenchmarkAsyncWindow(b *testing.B) {
	ratelimit.Benchmark(b, func() ratelimit.RateLimiter {
		return ratelimit.NewAsyncWindow(time.Millisecond, 10)
	})
}
