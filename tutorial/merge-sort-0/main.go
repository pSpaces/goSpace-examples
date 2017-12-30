package main

import (
	"fmt"

	. "github.com/pspaces/gospace"
)

//N defines the number of "splitter" workers
const N = 2

//M defines the number of "merger" workers
const M = 2

func main() {

	space := NewSpace("space")

	// We place the array to be sorted in the tuple space
	a := []int{7, 6, 5, 4, 3, 2, 1}
	space.Put("sort", a, len(a))

	// We add a lock for coordinating the merger workers
	space.Put("lock")

	// We launch all workers
	for i := 0; i < N; i++ {
		go splitter(&space, i)
	}
	for i := 0; i < M; i++ {
		go merger(&space, i)
	}

	// Here we wait for our result
	t, _ := space.Query("result", &a)
	fmt.Printf("RESULT: %v\n", (t.GetFieldAt(1)).([]int))

}

func splitter(space *Space, me int) {
	var a []int
	var resultLength int
	for {
		t, _ := space.Get("sort", &a, &resultLength)
		a = (t.GetFieldAt(1)).([]int)
		resultLength = (t.GetFieldAt(2)).(int)
		fmt.Printf("Splitter %d got %v...\n", me, a)
		// This should not happen
		if len(a) == 0 {
			continue
		}
		if len(a) == 1 {
			space.Put("sorted", a, 1, resultLength)
		} else {
			space.Put("sort", a[0:len(a)/2], resultLength)
			space.Put("sort", a[len(a)/2:len(a)], resultLength)
		}
	}
}

func merger(space *Space, me int) {
	var a []int
	var b []int
	var taskLength int
	var resultLength int
	for {
		// We use a lock to avoid deadlocks due to mutually waiting merger workers
		space.Get("lock")
		t, _ := space.Get("sorted", &a, &taskLength, &resultLength)
		a = (t.GetFieldAt(1)).([]int)
		taskLength = (t.GetFieldAt(2)).(int)
		resultLength = (t.GetFieldAt(3)).(int)
		fmt.Printf("Merger %d got %v...\n", me, a)
		if taskLength == resultLength {
			space.Put("result", a)
			space.Put("lock")
			break
		} else {
			t, _ := space.Get("sorted", &b, &taskLength, &resultLength)
			b = (t.GetFieldAt(1)).([]int)
			taskLength = (t.GetFieldAt(2)).(int)
			resultLength = (t.GetFieldAt(3)).(int)
			fmt.Printf("Merger %d got %v...\n", me, b)
			space.Put("lock")

			// Standard merge of two ordered vectors a and b
			c := merge(a, b)
			space.Put("sorted", c, len(c), resultLength)
		}
	}
}
func merge(a []int, b []int) []int {
	i := 0
	j := 0
	k := 0
	c := make([]int, len(a)+len(b))
	for {
		if i == len(a) {
			for ; j < len(b); j++ {
				c[k] = b[j]
				k++
			}
			break
		}
		if j == len(b) {
			for ; i < len(a); i++ {
				c[k] = a[i]
				k++
			}
			break
		}
		if a[i] <= b[j] {
			c[k] = a[i]
			i++
		} else {
			c[k] = b[j]
			j++
		}
		k++
	}
	return c
}
