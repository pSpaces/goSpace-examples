// This is a simple producer/consumer system

package main

import (
	"fmt"

	. "github.com/pspaces/gospace"
)

// PRODUCERS defines the number of producers
const PRODUCERS = 2

// NTASKS defines the number of tasks each producer generates
const NTASKS = 2

//CONSUMERS defines the number of consumers
const CONSUMERS = 2

func main() {

	bag := NewSpace("bag")

	for i := 0; i < PRODUCERS; i++ {
		go producer(&bag, i, NTASKS)
	}

	for i := 0; i < CONSUMERS; i++ {
		go consumer(&bag, i)
	}

	bag.Query("done")

}

func producer(bag *Space, me int, ntasks int) {
	for i := 0; i < ntasks; i++ {
		bag.Put("task", me, i)
		fmt.Printf("Producer %d put task %d in the bag...\n", me, i)
	}
}

func consumer(bag *Space, me int) {
	var task int
	var producer int
	for {
		bag.Get("task", &producer, &task)
		fmt.Printf("Consumer %d got task %d from producer %d.\n", me, task, producer)
	}
}
