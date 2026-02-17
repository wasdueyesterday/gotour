package gotour

import (
	"sync"
	"time"
	"math"
)

type userBucket struct {
	tokens float64
	lastRefill time.Time
}

type RateLimiter struct {
	// your fields
	rate float64
	capacity float64
	userBuckets map[string]*userBucket
	mu sync.RWMutex
}

func NewRateLimiter(req int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		// init
		rate: float64(req) / window.Seconds(),
		capacity: float64(req), 
		userBuckets: make(map[string]*userBucket),
	}
}

func (rl *RateLimiter) Allow(userID string) bool {
	// implement
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.userBuckets[userID]
	if !ok {
		b = &userBucket{
			tokens: rl.capacity,
			lastRefill: time.Now(),
		}
		rl.userBuckets[userID] = b
	}

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	if elapsed < 0 {
		elapsed = 0
	}
	newTokens := elapsed * rl.rate 
	b.tokens = math.Min(rl.capacity, b.tokens+newTokens)
	b.lastRefill = now

	if b.tokens >= 1 {
		b.tokens -= 1
		return true
	}
	return false
}

// Optional bonus
func (rl *RateLimiter) TokensFor(userID string) int {
	// return current tokens available for this user
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	b, ok := rl.userBuckets[userID]
	if !ok {
		return int(rl.capacity)
	}

	// Calculate current token without actually use one
	elapsed := time.Since(b.lastRefill).Seconds()
	currentTokens := b.tokens + (elapsed * rl.rate)

	return int(math.Min(rl.capacity, currentTokens))
}

