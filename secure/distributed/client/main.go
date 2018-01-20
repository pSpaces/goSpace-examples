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

	uri := strings.Join([]string{"tcp://", host, port, "/", name}, "")

	_, config := certificate.GenerateCertConfigs()

	spc := NewRemoteSpace(uri, config)

	// Put a message in the space.
	spc.Put("Hello, Alice!")
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
