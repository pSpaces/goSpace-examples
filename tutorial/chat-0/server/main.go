package main

import (
	"flag"
	"fmt"

	. "github.com/pspaces/gospace"
)

func main() {

	// Extract local host and port number from the command line
	host, port, space := args()

	// Create the chat space
	uri := "tcp://" + host + ":" + port + "/" + space
	fmt.Printf("Setting up chat space %s\n", uri)
	chat := NewSpace(uri)

	var who string
	var message string

	for {
		// Get and display chat messages
		t, _ := chat.Get(&who, &message)
		who = (t.GetFieldAt(0)).(string)
		message = (t.GetFieldAt(1)).(string)
		fmt.Printf("%s: %s \n", who, message)
	}

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
