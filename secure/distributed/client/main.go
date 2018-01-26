package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/pspaces/goSpace-examples/secure/certificate"
	. "github.com/pspaces/gospace"
)

func main() {
	host, port := args()
	if host == "" {
		return
	}

	name := "space"

	uri := strings.Join([]string{"tcp://", host, port, "/", name}, "")

	_, config := certificate.GenerateCertConfigs()

	spc := NewRemoteSpace(uri, config)

	// Put a message in the space.
	t, b := spc.Put("Hello, Alice!")

	fmt.Println("Received: ", t, b)
	println("Message put")

	time.Sleep(1 * time.Second)

	println("Getting message")

	var message string
	spc.Get(&message)

	println("Message get")

	fmt.Printf("%s\n", message)
}

func args() (host string, port string) {
	flag.Parse()

	argn := flag.NArg()
	if argn > 2 {
		fmt.Printf("Usage of %s: [address] [port]\n", "bob")
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
