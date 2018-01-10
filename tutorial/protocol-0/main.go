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
	"math/rand"
	"os"
	"strconv"
	"strings"

	. "github.com/pspaces/gospace"
)

func main() {

	// This is the server lobby
	lobby := NewSpace("tcp://localhost:31415/lobby")

	// create identity database with some users
	google := NewSpace("tcp://localhost:31416/google")
	google.Put("Alice", "1234")
	google.Put("Bob", "1234")

	// launch welcome
	go idProvider(&google)
	go welcome(&lobby, &google)

	// simulate a user
	for {
		// Read username
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		fmt.Print("Enter password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)
		client(&lobby, username, password)
	}

}

func welcome(lobby *Space, idProvider *Space) {
	var s string
	var username string
	sessionKey := 31417

	for {
		// get request for a new session of the login protocol
		t, _ := lobby.Get("login", &s)
		username = (t.GetFieldAt(1)).(string)
		fmt.Printf("Login request from %s...\n", username)
		sessionURL := "tcp://localhost:" + strconv.Itoa(sessionKey) + "/session" + strconv.Itoa(sessionKey)
		fmt.Printf("Creating session space %s...\n", sessionURL)
		session := NewSpace(sessionURL)
		session.QueryP("whatever")
		lobby.Put("session", username, sessionURL, "Client2Server", "Server2Client")
		go server(sessionURL, "Server2Client", "Client2Server", "Server2IdProvider", "IdProvider2Server")
		idProvider.Put("session", sessionURL, "IdProvider2Server", "Server2IdProvider")
		sessionKey++
		session.Get("go")
	}
}

func client(server *Space, username string, password string) {
	var s string
	var i int
	var sessionURL string
	var toServer string
	var fromServer string

	// request session of the protocol
	server.Put("login", username)
	t, _ := server.Get("session", username, &s, &s, &s)
	sessionURL = (t.GetFieldAt(2)).(string)
	toServer = (t.GetFieldAt(3)).(string)
	fromServer = (t.GetFieldAt(4)).(string)

	// run session of the protocol
	session := NewRemoteSpace(sessionURL)
	fmt.Printf("Starting client on session %s...\n", sessionURL)

	// The protocol
	session.Put("crap")
	fmt.Printf("Client: sending %s, %s, %s\n", toServer, username, password)
	session.Put(toServer, username, password)

	fmt.Printf("Client: waiting for the branch...\n")
	t, _ = session.Get(fromServer, &s)
	branch := (t.GetFieldAt(i)).(string)
	if branch == "true" {
		for {
			if enoughData() {
				session.Put(toServer, "continue")
				fmt.Printf("Client: waiting for data from the server...\n")
				t, _ = session.Get(fromServer, &i)
			} else {
				session.Put(toServer, "break")
				break
			}
		}
	} else {
		fmt.Printf("Client: waiting for bad news...\n")
		session.Get(fromServer, "ko")
	}
}

func enoughData() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("More? (yes/no)")
	reply, _ := reader.ReadString('\n')
	reply = strings.TrimSpace(reply)
	if reply == "yes" {
		return true
	}
	return false
}

func server(sessionURL string, toClient string, fromClient string, toIdProvider string, fromIdProvider string) {
	var s string
	var username string
	var password string
	var response string

	fmt.Printf("Starting server on session %s...\n", sessionURL)
	session := NewRemoteSpace(sessionURL)

	// The protocol
	fmt.Printf("Server: waiting for user credentials...\n")
	t, _ := session.Get(fromClient, &s, &s)
	username = (t.GetFieldAt(1)).(string)
	password = (t.GetFieldAt(2)).(string)
	fmt.Printf("Server got credentials ()%s,%s) on session %s...\n", sessionURL, username, password)
	session.Put(toIdProvider, username, password)
	fmt.Printf("Server: waiting for response from ID provider...\n")
	session.Get(fromIdProvider, &response)
	fmt.Printf("Server got response from id provider on session %s...\n", sessionURL)
	if response == "ok" {
		session.Put(toClient, "then")
		for {
			fmt.Printf("Server: waiting for the client to decide on continue/break...\n")
			session.Get(fromClient, &s)
			if (t.GetFieldAt(1)).(string) == "continue" {
				fmt.Printf("Server: waiting for user to ask for data...\n")
				_, _ = session.Get(fromClient, "getData")
				session.Put(toClient, rand.Intn(10))
			} else {
				break
			}
		}
	} else {
		session.Put(toClient, "else")
		session.Put(toClient, "ko")
	}
}

func idProvider(users *Space) {
	var s string
	var sessionURL string
	var toServer string
	var fromServer string

	for {
		t, _ := users.Get("session", &s, &s, &s)
		sessionURL = (t.GetFieldAt(1)).(string)
		toServer = (t.GetFieldAt(2)).(string)
		fromServer = (t.GetFieldAt(3)).(string)
		fmt.Printf("Starting id provider session %s...\n", sessionURL)
		go idProviderSession(users, sessionURL, toServer, fromServer)
	}
}

func idProviderSession(users *Space, sessionURL string, toServer string, fromServer string) {
	var s string
	var username string
	var password string
	session := NewRemoteSpace(sessionURL)

	// The protocol
	fmt.Printf("ID provider: waiting for user credentials...\n")
	t, _ := session.Get(fromServer, &s, &s)
	username = (t.GetFieldAt(1)).(string)
	password = (t.GetFieldAt(1)).(string)
	fmt.Printf("Id provider got credentials in session %s...\n", sessionURL)
	_, err := users.QueryP(username, password)
	if err == nil {
		session.Put(toServer, "ok")
	} else {
		session.Put(toServer, "ko")
	}
}
