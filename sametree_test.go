package gotour

import (
	"fmt"
	"testing"

	"golang.org/x/tour/tree"
)

func TestSametree(t *testing.T) {
	r := Same(tree.New(1), tree.New(1))
	fmt.Printf("r = %v\n", r)
	
	r2 := Same(tree.New(1), tree.New(2))
	fmt.Printf("r2 = %v\n", r2)
}
