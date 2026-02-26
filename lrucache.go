package gotour

import (
	"fmt"
	"sync"
)

type LRUNode struct {
	key, value int
	prev, next  *LRUNode // doubly linked lists
}

type LRUCache struct {
    capacity int
	
	mu sync.Mutex
	cache map[int]*LRUNode 

	// these are sentinal nodes that just marks the lru and mru
	// they are not data nodes!
	lru, mru *LRUNode
}


func NewLRUCache(capacity int) *LRUCache {
	left := &LRUNode{key: 0, value: 0}
	right := &LRUNode{key: 0, value: 0}

	left.next = right
	right.prev = left

    return &LRUCache{
		capacity: capacity,
		cache: make(map[int]*LRUNode),
		lru: left,
		mru: right,
	}
}


func (t *LRUCache) Get(key int) int {
	t.mu.Lock()
	defer t.mu.Unlock()

    v, exist := t.cache[key]
	if !exist {
		return -1
	} 

	// this pointer check to avoid unnecessary ops
	if v.next == t.mru {
		fmt.Println("Already before the mru sentinal!")
		return v.value
	}

	// removing from list first fix up the pointer
	// then insert into the map
	// this ordering is better preserved
	// remove v from the middle
	t.remove(v)
	
	// move v before mru
	t.insert(v)

	return v.value
}

// helper func to remove any node in the list
func (t *LRUCache) remove(node *LRUNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// helper func to always insert before the mru sentinel node
func (t *LRUCache) insert(node *LRUNode) {
	p := t.mru.prev
	p.next = node
	node.prev = p
	node.next = t.mru
	t.mru.prev = node
}

func (t *LRUCache) Put(key int, value int)  {
    t.mu.Lock()
	defer t.mu.Unlock()

	if v, exist := t.cache[key]; exist {
		// update the value
		v.value = value
		if v.next == t.mru {
			fmt.Println("Already next to mru!")
			return
		}
		t.remove(v)
		t.insert(v)
		return
	}

	// create the node
	nd := &LRUNode{
		key: key,
		value: value,
	}
	// check capacity before Put
	if len(t.cache) == t.capacity {
		// time to evict, the actual lru data node
		evict := t.lru.next

		// delete from cache
		delete(t.cache, evict.key)

		// delete from list
		t.remove(evict)
	
	}
	t.insert(nd)
	t.cache[key] = nd
}

