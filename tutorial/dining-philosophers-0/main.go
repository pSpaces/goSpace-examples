// This is a model of the classic problem of the dining philosophers.
// The solution is based on using tickets to limit concurrency.

package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
)

// N defines the number of philosophers.
const N = 10

func main() {
	board := NewSpace("board")

	go waiter(&board)

	for i := 0; i < N; i++ {
		go philosopher(&board, i)
	}

	board.Query("done")
}

// waiter prepares the board with forks.
func waiter(board *Space) {
	fmt.Printf("Waiter putting forks on the table...\n")

	for i := 0; i < N; i++ {
		board.Put("fork", i)
		fmt.Printf("Waiter put fork %d on the table.\n", i)
	}

	fmt.Printf("Waiter done.\n")
}

// philosopher defines the behaviour of a philosopher.
func philosopher(board *Space, me int) {
	// We define variables to identify the left and right forks.
	left := me
	right := (me + 1) % N

	// The philosopher enters his endless life cycle.
	for {
		// Wait until the left fork is ready (get the corresponding tuple).
		board.Get("fork", left)
		fmt.Printf("Philosopher %d got left fork\n", me)

		// Wait until the right fork is ready (get the corresponding tuple).
		board.Get("fork", right)
		fmt.Printf("Philosopher %d got right fork\n", me)

		// Lunch time.
		fmt.Printf("Philosopher %d is eating...\n", me)

		// Return the forks (put the corresponding tuples).
		board.Put("fork", left)
		board.Put("fork", right)
		fmt.Printf("Philosopher %d put both forks on the table\n", me)
	}
}
