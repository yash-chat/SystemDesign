package main

import (
	"sync"
	"time"
)

// SlidingLog implements a sliding log rate limiter
type SlidingLog struct {
	limit      int           // max requests allowed in window
	windowSize time.Duration // time window (e.g. 1 second)
	logs       []time.Time   // timestamps of recent requests
	mu         sync.Mutex
}

// NewSlidingLog creates a new sliding log rate limiter
func NewSlidingLog(limit int, windowSize time.Duration) *SlidingLog {
	return &SlidingLog{
		limit:      limit,
		windowSize: windowSize,
		logs:       make([]time.Time, 0),
	}
}

// Allow checks if a new request is allowed
func (sl *SlidingLog) Allow() bool {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sl.windowSize)

	// Step 1: Remove old timestamps outside the window
	validLogs := make([]time.Time, 0, len(sl.logs))
	for _, t := range sl.logs {
		if t.After(cutoff) {
			validLogs = append(validLogs, t)
		}
	}
	sl.logs = validLogs

	// Step 2: Check if limit exceeded
	if len(sl.logs) < sl.limit {
		sl.logs = append(sl.logs, now)
		return true
	}
	return false
}

// GetLogCount returns number of requests in current window
func (sl *SlidingLog) GetLogCount() int {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return len(sl.logs)
}
