package ratelimiter

import (
	"testing"
	"time"
)

func TestMuchHigherMaxRequests(t *testing.T) {
	numRequests := 1000
	limiter := NewQuantumRateLimiter("test", numRequests, time.Second)
	key := "foo"

	for i := 0; i < numRequests; i++ {
		// Should not be reached
		if limiter.Limit(key) {
			t.Errorf("N(%v) limit should not be reached.", i)
		}
	}

	if !limiter.Limit(key) {
		t.Errorf("N(%v) limit should be reached because it exceeds %v request per second.", numRequests+2, numRequests)
	}

}

func TestLimitNotReachedAfterRefill(t *testing.T) {
	numRequests := 1000
	limiter := NewQuantumRateLimiter("test", numRequests, time.Second)
	key := "foo"

	for i := 0; i < numRequests; i++ {
		if limiter.Limit(key) {
			t.Errorf("N(%v) limit should not be reached.", i)
		}
	}

	if !limiter.Limit(key) {
		t.Errorf("N(%v) limit should be reached because it exceeds %v request per second.", numRequests+2, numRequests)
	}
	// Waiting for the bucket to replenish
	time.Sleep(time.Second * 2)
	// Trying again, should be fine
	for i := 0; i < numRequests; i++ {
		if limiter.Limit(key) {
			t.Errorf("N(%v) limit should not be reached.", i)
		}
	}
	// Doing more, should all fail
	for i := 0; i < numRequests; i++ {
		if !limiter.Limit(key) {
			t.Errorf("N(%v) limit should not be reached.", i)
		}
	}
}
