package gotour

import (
	"sync"
)

type LFUNode struct {
	key, value int
	freq int // tracks the node's current frequency
	prev, next  *LFUNode
}

type List struct {
	lru, mru *LFUNode
	size int // for checking if a freq is empty
}

func NewList() *List {
	left := &LFUNode{}
	right := &LFUNode{}

	left.next = right
	right.prev = left

	return &List{
		lru: left,
		mru: right,
	}
}

type LFUCache struct {
    capacity int
	
		mu sync.Mutex
		cache map[int]*LFUNode 

		freqMap map[int]*List // each frequency vs a LRU list
		minFreq int // current lowest freq
}


func NewLFUCache(capacity int) *LFUCache {
    return &LFUCache{
			capacity: capacity,
			cache: make(map[int]*LFUNode),
			freqMap: make(map[int]*List),
		}
}

func (t *LFUCache) Get(key int) int {
		t.mu.Lock()
		defer t.mu.Unlock()

    v, exist := t.cache[key]
		if !exist {
			return -1
		} 

		t.promote(v)

		return v.value
}

func (t *LFUCache) removeFromList(node *LFUNode, lst *List) {
	node.prev.next = node.next
	node.next.prev = node.prev
	lst.size--
}

func (t *LFUCache) insertToList(node *LFUNode, lst *List) {
	p := lst.mru.prev
	p.next = node
	node.prev = p
	node.next = lst.mru
	lst.mru.prev = node
}

// Need to promote the frequency with 2 steps
func (t *LFUCache) promote(n *LFUNode) {
	// remove from lower freqMap 
	if lst, exist := t.freqMap[n.freq]; exist {
		t.removeFromList(n, lst)

		// Handle special case where the list goes empty
		if lst.size == 0 && n.freq == t.minFreq {
			t.minFreq++
		}
	}
	
	// promote to higher freqMap
	n.freq++
	
	if _, exist := t.freqMap[n.freq]; !exist {
		t.freqMap[n.freq] = NewList()
	}
	t.insertToList(n, t.freqMap[n.freq])
}

func (t *LFUCache) Put(key int, value int)  {
	if t.capacity == 0 {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	if v, exist := t.cache[key]; exist {
		// update the value
		v.value = value
		t.promote(v)
		return
	} 
	// create the node
	nd := &LFUNode{
		key: key,
		value: value,
		freq: 1,
	}
	// check capacity before Put
	if len(t.cache) == t.capacity {
		// locate the one to evict
		if l, exist := t.freqMap[t.minFreq]; exist {
			evict := l.lru.next

			// delete from cache
			// prevent finding a node that is about to be destroyed
			delete(t.cache, evict.key)

			// remove from the list
			t.removeFromList(evict, l)
		}
	}
	t.minFreq = 1 // lowest rank just joined
	if _, ok := t.freqMap[t.minFreq]; !ok {
		t.freqMap[t.minFreq] = NewList()
	}
	// ensure the node is linked before it can be found
	t.insertToList(nd, t.freqMap[t.minFreq])
	t.cache[key] = nd
}
