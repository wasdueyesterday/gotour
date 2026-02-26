package gotour

import "testing"

func TestLRUCache_Case1(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 10)
	c.Put(2, 20)

	if c.Get(1) != 10 {
		t.Errorf("Should retrieve 10 for key 1")
	}

	c.Put(3, 30) // should evict key 2
	if c.Get(2) != -1 {
		t.Errorf("Should've evicted key 2 and get value -1")
	}

	if c.Get(3) != 30 {
		t.Errorf("Should get 30 for key 3")
	}

	c.Put(1, 100) // update key 1, now MRU, key 3 is LRU
	c.Put(4, 40)  // key 3 should be evicted

	if c.Get(1) != 100 {
		t.Errorf("Should get 100 for key 1 due to update")
	}

	if c.Get(3) != -1 {
		t.Errorf("Should've evicted key 3" )
	}
}

func BenchmarkLRUCachePut(b *testing.B) {
	c := NewLRUCache(1000)

	// ResetTimer excludes the setup time (constructor) from the result
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Put(i % 2000, i) // Use mod to trigger some eviction
	}
}

func BenchmarkLRUCacheGet(b *testing.B) {
	c := NewLRUCache(1000)
	for i := 0; i < 1000; i++ {
		c.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(i % 1000)
	}
}

