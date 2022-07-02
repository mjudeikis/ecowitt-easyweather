package ratelimiter

import (
	"sync"
	"time"

	"github.com/juju/ratelimit"
)

var _ Interface = &quantumRateLimiter{}

type Interface interface {
	// Limit returns true if the ratelimit for the given entity has been reached
	Limit(id string) bool
	// Cap total token count for this bucket: Header: RateLimit-Limit: X
	Cap() int64
	// Name of the bucket
	Name() string
	// Available tells how many tokens left. Header: Ratelimit-Remaining: X
	Available(id string) int64
	// Wait tell how to to rate limit reset. Header: Ratelimit-Reset: Xs
	Wait(id string) time.Duration
}

// QuantumRateLimiter .
type quantumRateLimiter struct {
	buckets map[string]*ratelimit.Bucket
	rate    int
	per     time.Duration
	mu      sync.RWMutex
	name    string
}

// NewQuantumRateLimiter .
func NewQuantumRateLimiter(name string, rate int, per time.Duration) *quantumRateLimiter {
	return &quantumRateLimiter{
		name:    name,
		rate:    rate,
		per:     per,
		buckets: make(map[string]*ratelimit.Bucket),
	}
}

func (r *quantumRateLimiter) getOrCreate(id string) *ratelimit.Bucket {
	r.mu.RLock()
	limiter, ok := r.buckets[id]
	r.mu.RUnlock()
	if ok {
		return limiter
	}
	limiter = ratelimit.NewBucketWithQuantum(r.per, int64(r.rate), int64(r.rate))
	r.mu.Lock()
	r.buckets[id] = limiter
	r.mu.Unlock()
	return limiter
}

// Limit returns true if the ratelimit for the given entity has been reached
func (r *quantumRateLimiter) Name() string {
	return r.name
}

func (r *quantumRateLimiter) Cap() int64 {
	return int64(r.rate)
}

// Available returns rate limit availability
func (r *quantumRateLimiter) Available(id string) int64 {
	return r.getOrCreate(id).Available()
}

// Limit returns true if the ratelimit for the given entity has been reached
func (r *quantumRateLimiter) Limit(id string) bool {
	return r.Wait(id) != 0
}

// Wait returns the time to wait until available
func (r *quantumRateLimiter) Wait(id string) time.Duration {
	return r.getOrCreate(id).Take(1)
}

// WaitMaxDuration returns the time to wait until available, but with a max
func (r *quantumRateLimiter) WaitMaxDuration(id string, max time.Duration) (time.Duration, bool) {
	return r.getOrCreate(id).TakeMaxDuration(1, max)
}
