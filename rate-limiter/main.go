package main

import (
	"fmt"
	"time"
)

// RateLimiter is a common interface for any rate-limiting algorithm
type RateLimiter interface {
	Allow() bool
}

// simulateRequests runs N requests against a rate limiter and logs the output
func simulateRequests(limiter RateLimiter, numRequests int, interval time.Duration) {
	for i := 1; i <= numRequests; i++ {
		allowed := limiter.Allow()
		timestamp := time.Now().Format("15:04:05.000")

		if allowed {
			fmt.Printf("[%s] ✅ Request %d allowed\n", timestamp, i)
		} else {
			fmt.Printf("[%s] ❌ Request %d rejected\n", timestamp, i)
		}

		time.Sleep(interval)
	}
}

func main() {
	// Uncomment one of these to choose your algorithm 👇

	// 1️⃣ Token Bucket (example params)
	// limiter := NewTokenBucket(5, 2) // capacity=5, refill rate=2 tokens/sec

	// 2️⃣ Leaky Bucket (example params)
	// limiter := NewLeakyBucket(5, 1) // capacity=5, leak rate=1 req/sec

	// 3️⃣ Fixed Window (example params)
	limiter := NewFixedWindow(5, 1*time.Second) // limit=5 requests per 1-second window

	// Run simulation: 15 requests, 200ms apart
	simulateRequests(limiter, 15, 200*time.Millisecond)
}
