package main

import (
	"flag"
	"fmt"
	. "github.com/pspaces/gospace"
	"strings"
)

func main() {
	host, port := args()
	if host == "" {
		return
	}

	name := "space"

	uri := strings.Join([]string{"tcp://", host, port, "/", name}, "")

	spc := NewSpace(uri)

	// Get a message from the space
	// via pattern matching.
	var message string
	spc.Get(&message)

	fmt.Printf("%s\n", message)
}

func args() (host string, port string) {
	flag.Parse()

	argn := flag.NArg()
	if argn > 2 {
		fmt.Printf("Usage of %s: [address] [port]\n", "alice")
		return
	}

	if argn >= 1 {
		host = flag.Arg(0)
	} else {
		host = "localhost"
	}

	if argn == 2 {
		port = strings.Join([]string{":", flag.Arg(1)}, "")
	}

	return host, port
}
