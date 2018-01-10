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
	"fmt"
	"strconv"

	. "github.com/pspaces/gospace"
)

func main() {

	// This is the server lobby
	lobby := NewSpace("tcp://localhost:31415/lobby")

	// create identity database with some users
	google := NewRemoteSpace("tcp://localhost:31416/google")

	welcome(&lobby, &google)

}

func welcome(lobby *Space, idProvider *Space) {
	var username string
	sessionKey := 31417

	for {
		// get request for a new session of the login protocol
		t, _ := lobby.Get("login", &username)
		username = (t.GetFieldAt(1)).(string)
		fmt.Printf("Login request from %s...\n", username)
		sessionURL := "tcp://localhost:" + strconv.Itoa(sessionKey) + "/session" + strconv.Itoa(sessionKey)
		fmt.Printf("Creating session space %s...\n", sessionURL)
		session := NewSpace(sessionURL)
		go server(&session, "Server2Client", "Client2Server", "Server2IdProvider", "IdProvider2Server")
		lobby.Put("session", username, sessionURL, "Client2Server", "Server2Client")
		idProvider.Put("session", sessionURL, "IdProvider2Server", "Server2IdProvider")
		sessionKey++
	}
}

func server(session *Space, toClient string, fromClient string, toIdProvider string, fromIdProvider string) {
	var username string
	var password string
	var response string
	var branch string

	sessionURL := "some session"
	fmt.Printf("Starting server on session %s...\n", sessionURL)
	//session := NewRemoteSpace(sessionURL)

	// The protocol
	session.Put("crap")
	fmt.Printf("Server: waiting for user credentials...\n")
	t, _ := session.Get(fromClient, &username, &password)
	username = (t.GetFieldAt(1)).(string)
	password = (t.GetFieldAt(2)).(string)
	fmt.Printf("Server got credentials ()%s,%s) on session %s...\n", sessionURL, username, password)
	session.Put(toIdProvider, username, password)
	fmt.Printf("Server: waiting for response from ID provider...\n")
	session.Get(fromIdProvider, &response)
	fmt.Printf("Server got response %s from id provider on session %s...\n", response, sessionURL)
	if response == "ok" {
		session.Put(toClient, "then")
		for {
			fmt.Printf("Server: waiting for the client to decide on continue/break...\n")
			t, _ := session.Get(fromClient, &branch)
			branch = (t.GetFieldAt(1)).(string)
			if (t.GetFieldAt(1)).(string) == "continue" {
				fmt.Printf("Server: waiting for user to ask for data...\n")
				_, _ = session.Get(fromClient, "getData")
				session.Put(toClient, "some data")
			} else {
				break
			}
		}
	} else {
		session.Put(toClient, "else")
		session.Put(toClient, "ko")
	}
}
