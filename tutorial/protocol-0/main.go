// This is a simple login protocol
// Client -> Server : (login,name,pwd)
// Server -> IdProvider : (login,mame,pwd)
// IdProvider -> Server : "ok"
// Server -> Client
// IdProvider -> Server: "ko"
// Server -> Client : "ko"

package main

import (
	"fmt"
	"strconv"

	. "github.com/pspaces/gospace"
)

func main() {

	server := NewSpace("tcp://localhost:31415/server")
	google := NewSpace("tcp://localhost:31416/google")
	google.Put("Alice", "1234")
	google.Put("Bob", "1234")

	go client(&server, "Alice", "1234")
	go client(&server, "Bob", "1234")
	go client(&server, "Charlie", "1234")
	go welcome(&server, &google)
	go idProvider(&google)

	for {
	}
}

func welcome(lounge *Space, idProvider *Space) {
	var username string
	sessionKey := 33333
	for {
		// get request for a new session of the login protocol
		lounge.Get("login", &username)
		fmt.Printf("Login request from %s...\n", username)
		sessionURL := "tcp://localhost:" + strconv.Itoa(sessionKey) + "/session" + strconv.Itoa(sessionKey)
		fmt.Printf("Creating session space %s...\n", sessionURL)
		session := NewSpace(sessionURL)
		session.Put("whatever")
		lounge.Put("session", username, sessionURL, "ClientServer")
		go server(sessionURL, "ClientServer", "ServerIdP")
		idProvider.Put("session", sessionURL, "ServerIdP")
		sessionKey++
	}
}

func client(server *Space, username string, password string) {
	var sessionURL string
	var reply string
	var channel string

	// request session of the protocol
	server.Put("login", username)
	server.Get("session", username, &sessionURL, &channel)

	// run session of the protocol
	session := NewRemoteSpace(sessionURL)
	session.Put(channel, username, password)
	session.Get(channel, &reply)

	fmt.Printf("%s: got %s\n", username, reply)
}

func server(sessionURL string, clientChannel string, idProviderChannel string) {
	var username string
	var password string
	var reply string

	fmt.Printf("Starting server on session %s...\n", sessionURL)
	session := NewRemoteSpace(sessionURL)

	session.Get(clientChannel, &username, &password)
	fmt.Printf("Server got credentials on session %s...\n", sessionURL)
	session.Put(idProviderChannel, username, password)
	session.Get(idProviderChannel, &reply)
	fmt.Printf("Server got reply from id provider on session %s...\n", sessionURL)
	session.Put(clientChannel, reply)
}

func idProvider(users *Space) {
	var sessionURL string
	var channel string
	for {
		users.Get("session", &sessionURL, &channel)
		fmt.Printf("Starting id provider session %s...\n", sessionURL)
		go idProviderSession(users, sessionURL, channel)
	}
}

func idProviderSession(users *Space, sessionURL string, channel string) {
	var username string
	var password string
	session := NewRemoteSpace(sessionURL)
	session.Get(channel, &username, &password)
	fmt.Printf("Id provider got credentials in session %s...\n", sessionURL)
	_, err := users.QueryP(username, password)
	if err == nil {
		session.Put(channel, "ok")
	} else {
		session.Put(channel, "ko")
	}
}
