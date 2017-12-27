package main

import (
	"fmt"

	. "github.com/pspaces/gospace"
)

func main() {
	inbox := NewSpace("space")

	// Put a message into the space.
	inbox.Put("Hello world!")

	// Get a message from the space
	// via pattern matching.
	var message string
	inbox.Get(&message)

	fmt.Println(message)
}
