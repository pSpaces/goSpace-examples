package main

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/pspaces/goSpace-examples/secure/certificate"
	. "github.com/pspaces/gospace"
	"github.com/pspaces/gospace/shared"
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
	config, _ := certificate.GenerateCertConfigs()

	// Setup the space with the URI and config.
	spc := NewSpace(uri, config)

	// The is necessary for the current pSpace implementation.
	// TODO: remove this once merged into aggregation branch.
	shared.CreateTypeField(reflect.TypeOf("abc"))

	// Put a message in the space.
	t, _ := spc.Put("Hello, Alice!")
	fmt.Println("Put tuple securely into space: ", t)

	// Get a number from the space that doesn't exists.
	var number int
	spc.Get(&number)

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
