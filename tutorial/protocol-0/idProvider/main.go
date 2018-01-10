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

	. "github.com/pspaces/gospace"
)

func main() {

	// create identity database with some users
	google := NewSpace("tcp://localhost:31416/google")
	google.Put("Alice", "1234")
	google.Put("Bob", "1234")

	// launch id provider
	idProvider(&google)

}

func idProvider(users *Space) {
	var sessionURL string
	var toServer string
	var fromServer string

	for {
		t, _ := users.Get("session", &sessionURL, &toServer, &fromServer)
		sessionURL = (t.GetFieldAt(1)).(string)
		toServer = (t.GetFieldAt(2)).(string)
		fromServer = (t.GetFieldAt(3)).(string)
		fmt.Printf("Starting id provider session %s...\n", sessionURL)
		go idProviderSession(users, sessionURL, toServer, fromServer)
	}
}

func idProviderSession(users *Space, sessionURL string, toServer string, fromServer string) {
	var username string
	var password string
	session := NewRemoteSpace(sessionURL)

	// The protocol
	fmt.Printf("ID provider: waiting for user credentials...\n")
	t, _ := session.Get(fromServer, &username, &password)
	username = (t.GetFieldAt(1)).(string)
	password = (t.GetFieldAt(2)).(string)
	fmt.Printf("Id provider got credentials %v in session %s...\n", t, sessionURL)
	_, err := users.QueryP(username, password)
	if err == nil {
		session.Put(toServer, "ok")
	} else {
		session.Put(toServer, "ko")
	}
}
