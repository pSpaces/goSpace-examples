package main

import (
	"flag"
	"fmt"
	"strings"

	. "github.com/pspaces/gospace"
)

func main() {

	host, port := args()
	uri := strings.Join([]string{"tcp://", host, ":", port, "/space"}, "")
	space := NewRemoteSpace(uri)

	go splitter(&space, 0)

	for {
	}

}

func splitter(space *Space, me int) {
	var a []int
	var resultLength int
	for {
		_, err := space.Get("sort", &a, &resultLength)
		if err != nil {
			fmt.Println("Error!")
			return
		}
		fmt.Printf("Splitter %d got %v...\n", me, a)
		// This should not happen
		if len(a) == 0 {
			fmt.Printf("Splitter %d throwing away garbage...\n", me)
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

func args() (host string, port string) {

	// default values
	port = "31145"
	host = "localhost"

	flag.Parse()
	argn := flag.NArg()

	if argn > 2 {
		fmt.Printf("Too many arguments\nUsage: [host] [port]\n")
		return
	}

	if argn >= 1 {
		host = flag.Arg(0)
	}

	if argn >= 2 {
		port = flag.Arg(1)
	}

	return host, port
}
