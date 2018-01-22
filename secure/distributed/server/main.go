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

	config, _ := certificate.GenerateCertConfigs()

	name := "space"

	uri := strings.Join([]string{"tcp://", host, port, "/", name}, "")

	spc := NewSpace(uri, config)

	// Get a message from the space
	// via pattern matching.
	// var message string
	// spc.Get(&message)

	// fmt.Printf("%s\n", message)

	fmt.Printf("%v\n", spc)
	time.Sleep(20 * time.Minute)
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
