package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/pspaces/goSpace-examples/secure/certificate"
	. "github.com/pspaces/gospace"
)

func main() {
	host, port := args()
	if host == "" {
		return
	}

	name := "space"

	// Create URI.
	uri := strings.Join([]string{"tcp://", host, port, "/", name}, "")

	// Create config for authentication.
	_, config := certificate.GenerateCertConfigs()

	// Setup the remote space with the URI and config.
	spc := NewRemoteSpace(uri, config)

	// Get message from the remote space.
	var message string
	spc.Get(&message)

	// Print message received.
	fmt.Println("Message received through a secure channel: ", message)
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
