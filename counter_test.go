package gotour

import (
	"sync"
	"testing"
)

func TestSafeCounter(t *testing.T) {
	c := SafeCounter{v: make(map[string]int)}
	key := "testkey"
	iter := 1000

	var wg sync.WaitGroup	
	wg.Add(iter)

	for i := 0; i < iter; i++ {
		go func() {
			defer wg.Done()
			c.Inc(key)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check the final value of the counter
	if c.Value(key) != iter {
		t.Errorf("Expected counter value to be %d, but got %d", iter, c.Value(key))
	}
}