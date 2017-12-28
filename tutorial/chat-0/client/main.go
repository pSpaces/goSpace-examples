package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	. "github.com/pspaces/gospace"
)

func main() {

	// Extract hostname, port and space id from the arguments
	host, port, space := args()

	// Connect to the chat space
	uri := "tcp://" + host + ":" + port + "/" + space
	fmt.Printf("Connecting to chat space %s\n", uri)
	chat := NewRemoteSpace(uri)

	// Read name from the console
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pick a name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Keep sending whatever the user types
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		chat.Put(name, message)
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
