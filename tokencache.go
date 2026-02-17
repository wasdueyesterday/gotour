package gotour

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	Data  string
	ExpiresAt int64
}

type TokenCache struct {
	// your fields
	mu sync.RWMutex
	entries map[string]*CacheEntry
}

func NewTokenCache() *TokenCache {
	return &TokenCache{
		// init
		entries: make(map[string]*CacheEntry),
	}
}

func (c *TokenCache) Add(token, value string, ttl time.Duration) {
	// implement
	ce := &CacheEntry{
		Data: value,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[token]=ce
}

func (c *TokenCache) Get(token string) (string, bool) {
	// implement
	c.mu.RLock()
	defer c.mu.RUnlock()
	t, ok := c.entries[token]; if !ok {
		return "", false
	} 
	
	if t.ExpiresAt < time.Now().Unix() {
		fmt.Println("Token expired")
		return "", false
	} 

	return t.Data, true
	
}

func (c *TokenCache) deleteExpiredHelper() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().Unix()
	count := 0

	for tk, entry := range c.entries {
		if  entry.ExpiresAt < now {
			delete(c.entries, tk)
			count++
		}
	}
	return count
}

func (c *TokenCache) StartCleanup(ctx context.Context, interval time.Duration, wg *sync.WaitGroup) {
	// return stop func
	ticker := time.NewTicker(interval)
	wg.Add(1)
	go func() {
		defer ticker.Stop()
		defer wg.Done()
		for {
			select {
			case <- ticker.C:
				removed := c.deleteExpiredHelper()
				fmt.Printf("Deleted %v token expired on tiker\n", removed)
			case <- ctx.Done(): 
				fmt.Println("Quit on context cancel")
				return
			}
		}
	}()
}

