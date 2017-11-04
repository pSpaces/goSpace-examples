package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
)

func main() {
	fridge := NewSpace("fridge")

	// Add some stuff to the grocery list.
	fridge.Put("milk", 2)
	fridge.Put("butter", 3)

	// Retrieve one item via pattern matching.
	var item string
	var quantity int
	fridge.Get(&item, &quantity)

	// Print the item retrieved.
	fmt.Printf("%s: (%v, %v)\n", "Grocery item", item, quantity)
}
