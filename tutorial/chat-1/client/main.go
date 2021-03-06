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
	host, port, loungeId := args()

	// Connect to the chat space
	uri := "tcp://" + host + ":" + port + "/" + loungeId
	fmt.Printf("Connecting to chat space %s\n", uri)
	lounge := NewRemoteSpace(uri)

	// Read name from the console
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pick a name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// enter/create a rooms
	fmt.Print("Pick a room: ")
	var roomID string
	roomID, _ = reader.ReadString('\n')
	roomID = strings.TrimSpace(roomID)
	lounge.Put("enter", name, roomID)
	t, _ := lounge.Get("roomURI", name, roomID, &uri)
	uri = (t.GetFieldAt(3)).(string)
	fmt.Printf("Connecting to chat space %s\n", uri)
	room := NewRemoteSpace(uri)

	// Keep sending whatever the user types
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		room.Put(name, message)
	}
}

func args() (host string, port string, space string) {

	// default values
	host = "localhost"
	port = "31415"
	space = "lounge"

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
