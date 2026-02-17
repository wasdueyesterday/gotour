package gotour

import (
	"sync"
	"testing"
	"time"
	"fmt"
)


func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(1, 5*time.Second) // 1 requests per 10 seconds

	// Simulate concurrent requests
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			user := fmt.Sprintf("user%d", id%3) // 3 users
			if limiter.Allow(user) {
				fmt.Printf("User %s GoGo at %v\n", user, time.Now().Format(time.Stamp))
			} else {
				fmt.Printf("User %s StopStop at %v\n", user, time.Now().Format(time.Stamp))
			}
			time.Sleep(2 * time.Second) // simulate work
		}(i)
	}
	wg.Wait()

	// Wait longer to see refill
	time.Sleep(3 * time.Second)
	fmt.Println("After refill:")
	if limiter.Allow("user0") {
		fmt.Println("user0 allowed after refill")
	}
}