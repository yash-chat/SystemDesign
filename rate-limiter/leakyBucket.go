package main

import (
	"sync"
	"time"
)

// LeakyBucket implements a rate limiter using the leaky bucket algorithm
type LeakyBucket struct {
	leakRate float64    // Rate at which requests leak out (requests/second)
	capacity float64    // Maximum capacity of the bucket
	level    float64    // Current level in the bucket
	lastLeak time.Time  // Last time we calculated leakage
	mu       sync.Mutex // Mutex for thread safety
}

// NewLeakyBucket creates a new LeakyBucket with specified capacity and leak rate
func NewLeakyBucket(capacity, leakRate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		lastLeak: time.Now(),
	}
}

// Allow checks if a request can be accepted into the bucket
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastLeak).Seconds()

	// Leak the bucket
	lb.level -= lb.leakRate * elapsed
	if lb.level < 0 {
		lb.level = 0
	}
	lb.lastLeak = now

	// Check if there's space for this request
	if lb.level < lb.capacity {
		lb.level += 1.0
		return true
	}
	return false
}

// GetLevel returns the current level of the bucket
func (lb *LeakyBucket) GetLevel() float64 {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.level
}
