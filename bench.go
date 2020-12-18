package ratelimit

import (
	"fmt"
	"sync"
	"testing"
)

func Benchmark(b *testing.B, provider func() RateLimiter) {
	// Run the bench mark for increasing number of concurrent goroutines.
	for _, m := range []int{1, 4, 16, 64, 256, 1024, 4096, 16384} {
		b.Run(fmt.Sprintf("Concurrency_%d", m), benchmark(provider, m))
	}
}

func benchmark(provider func() RateLimiter, m int) func(b *testing.B) {
	return func(b *testing.B) {
		limiter := provider()
		var wg sync.WaitGroup
		wg.Add(m)
		for i := 0; i < m; i++ {
			go func(n int) {
				for i := 0; i < n; i++ {
					limiter.Request(fmt.Sprint(i))
				}
				wg.Done()
			}(b.N)
		}
		wg.Wait()
	}
}
