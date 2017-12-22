// This is a simple example tuple space used to represent a to-do list on a fridge.

package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
)

func main() {
	// Creating a tuple.
	var tuple Tuple = CreateTuple("milk", 1)
	fmt.Print("We just created tuple")
	fmt.Println(tuple)

	fmt.Print("The fields of ")
	fmt.Print(tuple)
	fmt.Print(" are ")
	fmt.Print(tuple.GetFieldAt(0))
	fmt.Print(" and ")
	fmt.Println(tuple.GetFieldAt(1))

	// Creating a space.
	fridge := NewSpace("fridge")

	// Adding tuples.
	fridge.Put("coffee", 1)
	fridge.Put("clean kitchen")
	fridge.Put("butter", 2)
	fridge.Put("milk", 3)

	// Looking for a tuple.
	_, err1 := fridge.QueryP("clean kitchen")
	if err1 == nil {
		fmt.Println("We need to clean the kitchen")
	}

	// Removing a tuple.
	_, err2 := fridge.GetP("clean kitchen")
	if err2 == nil {
		fmt.Println("Cleaning...")
	}

	// Looking for a tuple with pattern matching.
	var numberOfBottles int
	_, err3 := fridge.QueryP("milk", &numberOfBottles)

	// Updating a tuple.
	if err3 == nil && numberOfBottles <= 10 {
		fmt.Println("We plan to buy milk, but not enough...")
		fridge.GetP("milk", &numberOfBottles)
		fridge.Put("milk", numberOfBottles+1)
	}

	var item string
	var quantity int
	groceryList, _ := fridge.QueryAll(&item, &quantity)
	fmt.Println("Items to buy: ")
	fmt.Println(groceryList)
}
