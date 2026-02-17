package gotour

import (
	"golang.org/x/tour/tree"
)


// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	defer close(ch)
	
	var walker func(*tree.Tree)
	walker = func(n *tree.Tree) {
		if n == nil {
			return
		}
		walker(n.Left)
		ch <- n.Value
		walker(n.Right)
	}
	
	walker(t)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	go Walk(t1, ch1)
	
	ch2 := make(chan int)
	go Walk(t2, ch2)

	for {
		v1, ok1 := <- ch1
		v2, ok2 := <- ch2
		
		if !ok1 && !ok2 {
			return true
		}
		if ok1 != ok2 || v1 != v2 {
			return false
		}
	}
}


