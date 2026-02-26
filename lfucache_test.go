package gotour

import "testing"

func TestLFUCache_Basic(t *testing.T) {
	lfu := NewLFUCache(2)

	lfu.Put(1, 1)
	lfu.Put(2, 2)

	f1 := lfu.cache[1].freq
	f2 := lfu.cache[2].freq
	if f1 != 1 || f2 != 1 || f1 != f2 { // freq == 1 for both keys
		t.Errorf("Expected both freq 1, got f1=%d, f2=%d", f1, f2)
	}

	lfu.Get(1) // promoted key1's freq to 2, key2 stays freq 1
	f11 := lfu.cache[1].freq
	f22 := lfu.cache[2].freq
	if f11 != 2 || f22 != 1{
		t.Errorf("Expected f11=2, f22=1, got f11=%d, f22=%d", f11, f22)
	}

	lfu.Put(3, 3) // evicted key2
	if lfu.Get(2) != -1 {
		t.Errorf("Expected key2 evicted but not")
	}
}

func BenchmarkLFUCachePut(b *testing.B) {
	c := NewLFUCache(1000)

	// ResetTimer excludes the setup time (constructor) from the result
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Put(i % 2000, i) // Use mod to trigger some eviction
	}
}

func BenchmarkLFUCacheGet(b *testing.B) {
	c := NewLFUCache(1000)
	for i := 0; i < 1000; i++ {
		c.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(i % 1000)
	}
}

func TestLFUCache_Hammering(t *testing.T) {
	lfu := NewLFUCache(10)

	for i := 0; i < 1000; i++ {
		key := i % 15 // Generate collisions and evictions
		lfu.Put(key, i)

		// Check integrity
		if !lfu.checkIntegrity() {
			t.Fatalf("Integrity check failed at iteration %d", i)
		}
	}
}


