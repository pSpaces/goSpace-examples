package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
	"time"
)

func main() {
	store := NewSpace("store")

	AddBooks(&store)

	go Cashier(&store)

	go Customer(&store)

	time.Sleep(2 * time.Second)
}

// AddBooks adds books to the store.
func AddBooks(store *Space) {
	book := "Of Mice and Men"
	store.Put(book, 200)
}

// Cashier handles the payment.
func Cashier(store *Space) {
	for {
		// Get the payment from the tuple space.
		var payment int
		var book string
		store.Get("Payment", &book, &payment)

		// Find the price of the book.
		var price int
		store.Query(book, &price)

		// Check if the priced paid is equal to what the book costs.
		if price == payment {
			fmt.Printf("Received payment of %d for the book \"%s\".\n", payment, book)
			// Remove the book from the store.
			store.Get(book, payment)
		}
	}
}

// Customer buys a book.
func Customer(store *Space) {
	// Search for the book and save the price.
	var price int
	book := "Of Mice and Men"
	store.Query(book, &price)
	fmt.Printf("Checked price for book \"%s\". The price is %d.\n", book, price)

	// Place payment for the book.
	store.Put("Payment", book, price)
	fmt.Printf("Placed payment for book \"%s\", at the price of %d.\n", book, price)
}
