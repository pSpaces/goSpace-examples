package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	. "github.com/pspaces/gospace"
)

func main() {

	// Extract local host and port number from the command line
	host, port, loungeID := args()

	// Create the chat space
	loungeURI := "tcp://" + host + ":" + port + "/" + loungeID
	fmt.Printf("Setting up lounge space %s...\n", loungeURI)
	lounge := NewSpace(loungeURI)

	go loungeWelcome(&lounge, host, port)

	lounge.Get("stop")
}

func loungeWelcome(lounge *Space, host string, port string) {

	var who string
	var roomID string
	var roomURI string
	// This maps room identifiers to port numbers
	var rooms map[string]int
	rooms = make(map[string]int)
	// chartPort will be used to ensure every chat space has a unique port number
	chatPort, _ := strconv.Atoi(port)
	chatPort++
	for {
		// Process room login requests
		t, _ := lounge.Get("enter", &who, &roomID)
		who = (t.GetFieldAt(1)).(string)
		roomID = (t.GetFieldAt(2)).(string)
		fmt.Printf("%s requesting to enter %s...\n", who, roomID)
		_, ok := rooms[roomID]
		if ok {
			roomURI = "tcp://" + host + ":" + strconv.Itoa(rooms[roomID]) + "/" + roomID
		} else {
			fmt.Printf("Creating room %s for %s...\n", roomID, who)
			rooms[roomID] = chatPort
			chatPort++
			roomURI = "tcp://" + host + ":" + strconv.Itoa(rooms[roomID]) + "/" + roomID
			fmt.Printf("Setting up chat space %s...\n", roomURI)
			room := NewSpace(roomURI)
			go show(&room, roomID)
		}
		fmt.Printf("Telling %s to go for room %s on uri %s", who, roomID, roomURI)
		lounge.Put("roomURI", who, roomID, roomURI)
	}
}

func show(room *Space, roomID string) {
	var who string
	var message string
	for {
		t, _ := room.Get(&who, &message)
		who = (t.GetFieldAt(0)).(string)
		message = (t.GetFieldAt(1)).(string)
		fmt.Printf("ROOM %s | %s: %s \n", roomID, who, message)
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
