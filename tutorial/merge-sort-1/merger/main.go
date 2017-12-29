package main

import (
	"flag"
	"fmt"
	"os"

	. "github.com/pspaces/gospace"
)

func main() {

	host, port := args()
	uri := "tcp://" + host + ":" + port + "/space"
	space := NewRemoteSpace(uri)

	go merger(&space, 0)

	for {
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

func args() (host string, port string) {

	// default values
	port = "31145"
	host = "localhost"

	flag.Parse()
	argn := flag.NArg()

	if argn > 2 {
		fmt.Println("Too many arguments")
		fmt.Println("Usage: go run main.go [host] [port]")
		os.Exit(-1)
	}

	if argn >= 1 {
		host = flag.Arg(0)
	}

	if argn >= 2 {
		port = flag.Arg(1)
	}

	return host, port
}
