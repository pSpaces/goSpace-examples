package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
)

func main() {
	spc := NewSpace("8080")

	// Put a message into the space.
	spc.Put("Hello, universe!")

	// Get a message from the space
	// via pattern matching.
	var message string
	spc.Get(&message)

	fmt.Println(message)
}
