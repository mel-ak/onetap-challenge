package usecases

import (
	"context"
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	rate       int
	interval   time.Duration
	tokens     int
	lastUpdate time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter with the specified rate and interval
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		interval:   interval,
		tokens:     rate,
		lastUpdate: time.Now(),
	}
}

// Wait blocks until a token is available or the context is cancelled
func (rl *RateLimiter) Wait(ctx context.Context) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Calculate time elapsed since last update
	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate)

	// Add new tokens based on elapsed time
	newTokens := int(elapsed.Seconds() * float64(rl.rate) / rl.interval.Seconds())
	if newTokens > 0 {
		rl.tokens = min(rl.rate, rl.tokens+newTokens)
		rl.lastUpdate = now
	}

	// If no tokens available, wait
	if rl.tokens <= 0 {
		waitTime := time.Duration(float64(rl.interval) / float64(rl.rate))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			rl.tokens = 1
			rl.lastUpdate = time.Now()
		}
	}

	rl.tokens--
	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
