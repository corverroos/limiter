package ratelimit

// RateLimiter is the interface implemented by a resource rate limiter.
type RateLimiter interface {
	// Request returns true if the request of the resource conforms
	// to the rate limit; it contributes to the rate. It returns false
	// if the request would exceed the rate limit; it doesn't
	// contributes to the rate.
	Request(resource string) bool
}
