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

	numPhilosophers, me, host, port := args()

	uri := strings.Join([]string{"tcp://", host, ":", port, "/board"}, "")
	board := NewRemoteSpace(uri)

	go philosopher(&board, me, numPhilosophers)

	board.Query("done")
}

// philospher defines the behaviour of a philosopher.
func philosopher(board *Space, me int, numPhilosophers int) {
	// We define variables to identify the left and right forks.
	left := me
	right := (me + 1) % numPhilosophers

	// The philosopher enters his endless life cycle.
	for {
		// Get a ticket.
		board.Get("ticket")

		// Wait until the left fork is ready (get the corresponding tuple).
		board.Get("fork", left)
		fmt.Printf("Philosopher %d got left fork\n", me)

		// Wait until the right fork is ready (get the corresponding tuple).
		board.Get("fork", right)
		fmt.Printf("Philosopher %d got right fork\n", me)

		// Lunch time.
		fmt.Printf("Philosopher %d is eating...\n", me)

		// Return the forks and the ticket (put the corresponding tuples).
		board.Put("fork", left)
		board.Put("fork", right)
		board.Put("ticket")
		fmt.Printf("Philosopher %d put both forks and a ticket on the table\n", me)
	}
}

func args() (numPhilosophers int, me int, host string, port string) {

	// default values
	numPhilosophers = 0
	me = 0
	port = "31145"
	host = "localhost"

	flag.Parse()
	argn := flag.NArg()

	if argn > 4 {
		fmt.Printf("Too many arguments\nUsage: [number of philosopers] [my id] [host] [port]\n")
		return
	}

	if argn >= 1 {
		numPhilosophers, _ = strconv.Atoi(flag.Arg(0))
	}

	if argn >= 2 {
		me, _ = strconv.Atoi(flag.Arg(1))
	}

	if argn >= 3 {
		host = flag.Arg(2)
	}

	if argn >= 4 {
		port = flag.Arg(3)
	}

	return numPhilosophers, me, host, port
}
