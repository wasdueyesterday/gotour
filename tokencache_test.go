package gotour

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTokenCache(t *testing.T) {
	cache := NewTokenCache()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	cache.StartCleanup(ctx, 2*time.Second, &wg)

	cache.Add("user123", "session-data", 5*time.Second)
	time.Sleep(3 * time.Second)
	val, ok := cache.Get("user123")
	fmt.Printf("After 3s: %v %v\n", val, ok) // expect value, true

	time.Sleep(3 * time.Second)
	val, ok = cache.Get("user123")
	fmt.Printf("After 6s: %v %v\n", val, ok) // expect "", false

	cancel()
	wg.Wait()
}