package main

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket implements a thread-safe token bucket rate limiter.
// It allows requests at a controlled rate, refilling tokens over time.
type TokenBucket struct {
	capacity float64    // Maximum number of tokens in the bucket
	fillRate float64    // Tokens added per second
	tokens   float64    // Current number of tokens
	lastFill time.Time  // Last time tokens were refilled
	mu       sync.Mutex // Mutex for thread safety
}

// NewTokenBucket creates a new TokenBucket with given capacity and fill rate.
func NewTokenBucket(capacity, fillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		fillRate: fillRate,
		tokens:   capacity, // Start full
		lastFill: time.Now(),
	}
}

// Allow checks if a request can proceed. Returns true if allowed, false otherwise.
// Refills tokens based on elapsed time since last check.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastFill).Seconds()

	// Refill tokens based on elapsed time and fill rate
	tb.tokens = min(tb.capacity, tb.tokens+elapsed*tb.fillRate)
	tb.lastFill = now

	if tb.tokens >= 1 {
		tb.tokens -= 1 // Consume a token
		return true
	}
	return false // Not enough tokens
}

// min returns the smaller of two float64 values.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Example usage: 10 tokens max, refill 1 token per second
	limiter := NewTokenBucket(10, 1)

	for i := 0; i < 20; i++ {
		if limiter.Allow() {
			fmt.Printf("[%s] Request allowed ✅\n", time.Now().Format("15:04:05.000"))
		} else {
			fmt.Printf("[%s] Request denied ❌ (no tokens left)\n", time.Now().Format("15:04:05.000"))
		}
		time.Sleep(300 * time.Millisecond)
	}
}
