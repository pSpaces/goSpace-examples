// user.(username,password) -> server.credentials(string,string) ;
// server.(username(credentials),password(credentials)) -> identityProvider.credentials(string,string) ;
// identityProvider.Put(check(credentials)) -> server.response(string) ;
// if response then
//     while notEnoughData()@client do {
//         client.("getData") -> server.t("getData")
//         server.(generateData()) ->client.data(int))
//     }
// else
//     server.("ko") -> user.response("ko")

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/pspaces/gospace"
)

func main() {

	// This is the server lobby
	lobby := NewRemoteSpace("tcp://localhost:31415/lobby")

	reader := bufio.NewReader(os.Stdin)
	for {
		// Read username
		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		fmt.Print("Enter password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)
		client(&lobby, username, password, reader)
	}

}

func client(server *Space, username string, password string, input *bufio.Reader) {
	var sessionURL string
	var toServer string
	var fromServer string
	var branch string
	var data string

	// request session of the protocol
	fmt.Printf("Client: sending login request...\n")
	server.Put("login", username)
	fmt.Printf("Client: waiting for session to be created...\n")
	t, _ := server.Get("session", username, &sessionURL, &toServer, &fromServer)
	sessionURL = (t.GetFieldAt(2)).(string)
	toServer = (t.GetFieldAt(3)).(string)
	fromServer = (t.GetFieldAt(4)).(string)

	// run session of the protocol
	session := NewRemoteSpace(sessionURL)
	fmt.Printf("Starting client on session %s...\n", sessionURL)

	// The protocol
	fmt.Printf("Client: sending %s, %s, %s\n", toServer, username, password)
	session.Put(toServer, username, password)

	fmt.Printf("Client: waiting for the branch...\n")
	t, _ = session.Get(fromServer, &branch)
	branch = (t.GetFieldAt(1)).(string)
	if branch == "then" {
		for {
			if enoughData(input) {
				session.Put(toServer, "continue")
				session.Put(toServer, "getData")
				fmt.Printf("Client: waiting for data from the server...\n")
				t, _ = session.Get(fromServer, &data)
				data = (t.GetFieldAt(1)).(string)
				fmt.Printf("Client: got %s\n", data)
			} else {
				session.Put(toServer, "break")
				break
			}
		}
	} else {
		fmt.Printf("Client: waiting for bad news...\n")
		session.Get(fromServer, "ko")
	}
	fmt.Printf("Client: done.\n")
}

func enoughData(input *bufio.Reader) bool {
	fmt.Print("More data?: ")
	reply, _ := input.ReadString('\n')
	reply = strings.TrimSpace(reply)
	if reply == "yes" {
		return true
	}
	return false
}
