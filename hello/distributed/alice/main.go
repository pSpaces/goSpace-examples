package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
)

func main() {
	spc := NewSpace("8080")

	// Get a message from the space
	// via pattern matching.
	var message string
	spc.Get(spc, &message)

	fmt.Println(message)
}
