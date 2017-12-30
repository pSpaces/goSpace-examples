// This is a simple example tuple space used to represent a to-do list on a fridge.

package main

import (
	"fmt"

	. "github.com/pspaces/gospace"
)

func main() {

	// Creating a space.
	fridge := NewSpace("fridge")

	go alice(&fridge)
	go bob(&fridge)
	go charlie(&fridge)

	fridge.Get("done")

}

func alice(fridge *Space) {
	fridge.Put("milk", 2)
	fridge.Put("butter", 1)
	fridge.Put("shop!")
}

func bob(fridge *Space) {
	var item string
	var quantity int
	fridge.Query("shop!")
	for {
		t, _ := fridge.Get(&item, &quantity)
		item = (t.GetFieldAt(0)).(string)
		quantity = (t.GetFieldAt(1)).(int)
		fmt.Printf("Bob: I am shopping %d items of %s...\n", quantity, item)
	}
}

func charlie(fridge *Space) {
	var item string
	var quantity int
	_, err := fridge.GetP("shop!")
	if err == nil {
		for {
			t, _ := fridge.Get(&item, &quantity)
			item = (t.GetFieldAt(0)).(string)
			quantity = (t.GetFieldAt(1)).(int)
			fmt.Printf("Charlie: I am shopping %d items of %s...\n", quantity, item)
		}
	} else {
		fmt.Printf("Charlie: I am just relaxing...\n")
	}
}
