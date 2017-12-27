// This is a model of the classic problem of the dining philosophers.
// The protocol uses tickets to limite concurrency and to avoid deadlocks.

package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	. "github.com/pspaces/gospace"
)

func main() {

	numPhilosophers, port := args()

	uri := strings.Join([]string{"tcp://localhost:", port, "/board"}, "")
	board := NewSpace(uri)

	go waiter(&board, numPhilosophers)

	board.Query("done")
}

// waiter prepares the board with forks and tickets.
func waiter(board *Space, numPhilosophers int) {
	fmt.Printf("Waiter putting forks on the table...\n")

	for i := 0; i < numPhilosophers; i++ {
		board.Put("fork", i)
		fmt.Printf("Waiter put fork %d on the table.\n", i)
	}

	fmt.Printf("Waiter putting tickets on the table...\n")

	for i := 0; i < numPhilosophers-1; i++ {
		board.Put("ticket")
	}

	fmt.Printf("Waiter done.\n")
}

func args() (numPhilosophers int, port string) {

	// default values
	numPhilosophers = 0
	port = "31145"

	flag.Parse()
	argn := flag.NArg()

	if argn > 2 {
		fmt.Printf("Too many arguments\nUsage: [number of philosopers] [port]\n")
		return
	}

	if argn >= 1 {
		numPhilosophers, _ = strconv.Atoi(flag.Arg(0))
	}

	if argn >= 2 {
		port = flag.Arg(1)
	}

	return numPhilosophers, port
}
