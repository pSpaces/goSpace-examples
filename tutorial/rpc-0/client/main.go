package main

import (
	"flag"
	"fmt"

	. "github.com/pspaces/gospace"
)

func main() {

	// Extract hostname, port and space id from the arguments
	host, port, space := args()

	// Connect to the chat space
	uri := "tcp://" + host + ":" + port + "/" + space
	fmt.Printf("Connecting to chat space %s\n", uri)
	server := NewRemoteSpace(uri)

	// Invoke foo(1,2,3) remotely
	fmt.Println("Invoking foo(1,2,3) on server...")
	server.Put("Alice1", "func", "foo")
	server.Put("Alice1", "args", 1, 2, 3)

	// Get the result
	var z int
	server.Get("Alice1", "result", &z)
	fmt.Printf("Server says foo(1,2,3) = %d \n", z)

	// Invoke bar("a","b") remotely
	server.Put("Alice2", "func", "bar")
	server.Put("Alice2", "args", "a", "b")

	// Get the result
	var c string
	server.Get("Alice2", "result", &c)
	fmt.Printf("Server says bar(\"a\",\"b\") = %s \n", c)

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
		return
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
