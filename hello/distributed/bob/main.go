package main

import (
	. "github.com/pspaces/gospace"
)

func main() {
	spc := NewRemoteSpace("8080")

	// Put a message in the space.
	spc.Put("Hello, Alice!")
}
