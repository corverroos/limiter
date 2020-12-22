package ratelimit

import (
	"fmt"
	"runtime"
	"testing"
)

func Benchmark(b *testing.B, provider func() RateLimiter) {
	// Run the benchmark for increasing number of concurrent goroutines.
	for _, m := range []int{1, 4, 16, 64, 256, 1024, 4096, 16384} {
		BenchmarkFunc(provider, m)(b)
	}
}

func BenchmarkFunc(provider func() RateLimiter, m int) func(b *testing.B){
	// Run the benchmark for increasing number of concurrent goroutines.
	return func(b *testing.B) {
		cpus := runtime.GOMAXPROCS(0)
		b.Run(fmt.Sprintf("Concurrency_%d", m), func(b *testing.B) {
			limiter := provider()
			b.SetParallelism(m/cpus)
			b.RunParallel(func(b *testing.PB) {
				for b.Next() {
					limiter.Request("")
				}
			})
		})
	}
}
