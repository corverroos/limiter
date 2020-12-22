# ratelimit

Package `ratelimit` provides a simple "resource based" rate limit API and multiple implementations.

The aim is to find a performant implementation to use as generic logging rate limiter for Luno microservices. We want to limit specific logs from spamming the infrastructure, identified either by unique message or source. Performance is more important than correctness in this use case since we expect high level of concurrency and slight over (or under) limiting is not an issue.

```
// RateLimiter is the interface implemented by a resource rate limiter.
type RateLimiter interface {
	// Request returns true if the request of the resource conforms
	// to the rate limit; it contributes to the rate. It returns false
	// if the request would exceed the rate limit; it doesn't
	// contributes to the rate.
	Request(resource string) bool
}
```

The package also provides performance benchmarking; number of operations per number of concurrent goroutines. The following are results of `ratelimit.Benchmark` (in ns/op) for different implementations with config aimed at `limit=100, period=10ms` executed on MacBook Pro Quad Core 16GB RAM: 

| Implementation | 1 | 4 | 16 | 64 | 256 | 1024 | 4096 | 16384 |
|---|--|--|--|--|--|--|--|--|
| github.com/corver/ratelimit.SyncMapWindow | 190 | 143 | 135 | 139 | 162 | 136 | 190 | 150 |
| github.com/corver/ratelimit.NaiveWindow | 223 | 187 | 240 | 235 | 251 | 358 | 361 | 395 |
| github.com/ellemouton/ratelimiter.TimerWindow | 59 | 63 | 81 | 111 | 116 | 122 | 121 | 121 |
| github.com/ellemouton/ratelimiter.ChannelWindow | 189 | 192 | 206 | 216 | 287 | 310 | 326 | 340 |

Please DM me link to your implementation for me to add your benchmark score.
