package gotour

import "testing"

func TestLruCache_Case1(t *testing.T) {
	c := NewLruCache(2)
	c.Put(1, 10)
	c.Put(2, 20)

	if c.Get(1) != 10 {
		t.Fatalf("Should retrieve 10 for key 1")
	}

	c.Put(3, 30) // should evict key 2
	if c.Get(2) != -1 {
		t.Fatalf("Should've evicted key 2 and get value -1")
	}

	if c.Get(3) != 30 {
		t.Fatalf("Should get 30 for key 3")
	}

	c.Put(1, 100) // update key 1, now MRU, key 3 is LRU
	c.Put(4, 40)  // key 3 should be evicted

	if c.Get(1) != 100 {
		t.Fatalf("Should get 100 for key 1 due to update")
	}

	if c.Get(3) != -1 {
		t.Fatalf("Should've evicted key 3" )
	}
}

