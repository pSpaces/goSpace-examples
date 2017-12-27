package main

import (
	"flag"
	"fmt"
	"strings"

	. "github.com/pspaces/gospace"
)

func main() {

	port := args()

	uri := strings.Join([]string{"tcp://localhost:", port, "/space"}, "")
	space := NewSpace(uri)

	// We place the array to be sorted in the tuple space
	a := []int{7, 6, 5, 4, 3, 2, 1}
	space.Put("sort", a, len(a))

	// We add a lock for coordinating the merger workers
	space.Put("lock")

	// Here we wait for our result
	space.Query("result", &a)
	fmt.Printf("RESULT: %v\n", a)

}

func args() (port string) {

	// default values
	port = "31145"

	flag.Parse()
	argn := flag.NArg()

	if argn > 1 {
		fmt.Printf("Too many arguments\nUsage: [port]\n")
		return
	}

	if argn >= 1 {
		port = flag.Arg(1)
	}

	return port
}
