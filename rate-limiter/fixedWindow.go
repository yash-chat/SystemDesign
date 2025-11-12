package main

import (
	"sync"
	"time"
)

// FixedWindow implements a fixed window rate limiter
type FixedWindow struct {
	limit       float64       // max requests allowed per window
	windowSize  time.Duration // duration of the window (e.g., 1 second)
	requests    float64       // number of requests in the current window
	windowStart time.Time     // start time of the current window
	mu          sync.Mutex
}

// NewFixedWindow creates a new fixed window rate limiter
func NewFixedWindow(limit float64, windowSize time.Duration) *FixedWindow {
	return &FixedWindow{
		limit:       limit,
		windowSize:  windowSize,
		windowStart: time.Now(),
	}
}

// Allow determines if a request is allowed under the rate limit
func (fw *FixedWindow) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()

	// Check if current window expired
	if now.Sub(fw.windowStart) >= fw.windowSize {
		// reset window
		fw.windowStart = now
		fw.requests = 0
	}

	// Check if within limit
	if fw.requests < fw.limit {
		fw.requests++
		return true
	}
	return false
}

// GetRequests returns the current request count in the window
func (fw *FixedWindow) GetRequests() float64 {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.requests
}
