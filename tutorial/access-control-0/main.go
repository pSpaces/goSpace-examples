// This is a simple access control system

package main

import (
	"fmt"

	. "github.com/pspaces/gospace"
)

func main() {

	// Space with the data to be protected
	space := NewSpace("tcp://localhost:31415/space")
	space.Put("inbox", "Alice", "Bob", "Hello")
	space.Put("board", "Alice", "Charlie", "Hello")

	// Space with access control rules
	rules := NewSpace("tcp://localhost:31416/rules")
	var s string
	// Rules for the inbox
	rules.Put("allow", "Alice", "Query", CreateTemplate("inbox", "Alice", &s, &s))
	rules.Put("allow", "Alice", "QueryP", CreateTemplate("inbox", "Alice", &s, &s))
	rules.Put("allow", "Bob", "Query", CreateTemplate("inbox", "Bob", &s, &s))
	rules.Put("allow", "Bob", "QueryP", CreateTemplate("inbox", "Bob", &s, &s))

	// Rules for the board
	rules.Put("allow", "Alice", "Query", CreateTemplate("board", "Alice", &s, &s))
	rules.Put("allow", "Alice", "Query", CreateTemplate("board", "Bob", &s, &s))
	rules.Put("allow", "Alice", "QueryP", CreateTemplate("board", "Alice", &s, &s))
	rules.Put("allow", "Alice", "QueryP", CreateTemplate("board", "Bob", &s, &s))
	rules.Put("allow", "Bob", "Query", CreateTemplate("board", "Alice", &s, &s))
	rules.Put("allow", "Bob", "Query", CreateTemplate("board", "Bob", &s, &s))
	rules.Put("allow", "Bob", "QueryP", CreateTemplate("board", "Alice", &s, &s))
	rules.Put("allow", "Bob", "QueryP", CreateTemplate("board", "Bob", &s, &s))

	// Space to interact with users
	lobby := NewSpace("tcp://localhost:31417/lobby")
	go server(&lobby, &rules, &space)

	user("Alice", &lobby)
	user("Bob", &lobby)
	user("Trend", &lobby)

	space.Query("stop")

}

func user(me string, lobby *Space) {
	var s string
	var t Tuple
	var err string

	lobby.Put(me, "QueryP", CreateTemplate("board", "Alice", &s, &s))
	t, _ = lobby.Get("reply", &s)
	if t.GetFieldAt(1) == "allowed" {
		t, _ = lobby.Get(me, "result", &t, &err)
		fmt.Printf("%s: I read %v on Alice's board\n", me, t.GetFieldAt(2))
	}

	lobby.Put(me, "QueryP", CreateTemplate("inbox", "Alice", &s, &s))
	t, _ = lobby.Get("reply", &s)
	if t.GetFieldAt(1) == "allowed" {
		t, _ = lobby.Get(me, "result", &t, &err)
		fmt.Printf("%s: I read %v on Alice's inbox\n", me, t.GetFieldAt(2))
	}
}

func server(lobby *Space, rules *Space, space *Space) {
	var subject string
	var action string
	var template Template
	var decision string
	for {
		t, _ := lobby.Get(&subject, &action, &template)
		subject = (t.GetFieldAt(0)).(string)
		action = (t.GetFieldAt(1)).(string)
		template = (t.GetFieldAt(2)).(Template)
		fmt.Printf("%s wants to do %s(%v)\n", subject, action, template)
		fmt.Printf("Checking rules...\n")
		decision = check(rules, subject, action, template)
		if decision != "permit" {
			fmt.Printf("Permission denied.\n")
			lobby.Put("reply", "denied")
			continue
		}
		fmt.Printf("Permission granted.\n")
		lobby.Put("reply", "allowed")
		fmt.Printf("Performing action...\n")
		switch action {
		case "Query":
			t, err := space.Query(template.Fields()...)
			lobby.Put(subject, "result", t, msg(err))
		case "QueryP":
			t, err := space.QueryP(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lobby.Put(subject, "result", t, msg(err))
		case "Get":
			t, err := space.Get(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lobby.Put(subject, "result", t, msg(err))
		case "GetP":
			t, err := space.Query(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lobby.Put(subject, "result", t, msg(err))
		case "Put":
			t, err := space.Query(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lobby.Put(subject, "result", t, msg(err))
		case "QueryAll":
			tl, err := space.QueryAll(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lobby.Put("result", tl, err)
		case "GetAll":
			tl, err := space.GetAll(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lobby.Put("result", tl, err)
		default:
			fmt.Printf("Providing result...\n")
			lobby.Put("result", nil, nil)
		}
	}
}

func msg(err error) string {
	if err == nil {
		return "ok"
	}
	return "ko"
}

func check(rules *Space, subject string, action string, template Template) string {
	var decision string
	decision = "deny"
	_, err := rules.QueryP(&decision, subject, action, template)
	if err == nil {
		decision = "permit"
	}
	return decision
}
