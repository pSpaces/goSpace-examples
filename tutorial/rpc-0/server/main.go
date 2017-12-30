package main

import (
	"flag"
	"fmt"
	"os"

	. "github.com/pspaces/gospace"
)

func main() {

	// Extract local host and port number from the command line
	host, port, space := args()

	// Create the chat space
	uri := "tcp://" + host + ":" + port + "/" + space
	fmt.Printf("Setting up RPC space %s\n", uri)
	mySpace := NewSpace(uri)

	var callID string
	var f string
	var x int
	var y int
	var z int
	var a string
	var b string

	for {
		// Get RPC request
		t, _ := mySpace.Get(&callID, "func", &f)
		callID = (t.GetFieldAt(0)).(string)
		f = (t.GetFieldAt(2)).(string)
		fmt.Printf("RPC %s received: f(%s", callID, f)
		switch f {
		case "foo":
			t, _ := mySpace.Get(callID, "args", &x, &y, &z)
			x = (t.GetFieldAt(2)).(int)
			y = (t.GetFieldAt(3)).(int)
			z = (t.GetFieldAt(4)).(int)
			fmt.Printf("%d,%d,%d)...\n", x, y, z)
			fmt.Println("Computing RPC...")
			result := foo(x, y, z)
			fmt.Println("Placing result...")
			mySpace.Put(callID, "result", result)

		case "bar":
			t, _ := mySpace.Get(callID, "args", &a, &b)
			a = (t.GetFieldAt(2)).(string)
			b = (t.GetFieldAt(3)).(string)
			fmt.Printf("%s,%s)...\n", a, b)
			fmt.Println("Computing RPC...")
			result := bar(a, b)
			fmt.Println("Placing result...")
			mySpace.Put(callID, "result", result)
		default:
			// ignore RPC for unknown functions
			fmt.Printf("...)...\n")
			fmt.Println("Ignoring request...")
			continue
		}
	}

}

func foo(x int, y int, z int) int {
	return x + y + z
}

func bar(a string, b string) string {
	return a + b + b + a
}

func args() (host string, port string, space string) {

	// default values
	host = "localhost"
	port = "31415"
	space = "chat"

	flag.Parse()
	argn := flag.NArg()

	if argn > 3 {
		fmt.Println("Too many arguments")
		fmt.Println("Usage: go run main.go [address] [port] [space]")
		os.Exit(-1)
	}

	if argn >= 1 {
		host = flag.Arg(0)
	}

	if argn >= 2 {
		port = flag.Arg(1)
	}

	if argn == 3 {
		space = flag.Arg(2)
	}

	return host, port, space
}
