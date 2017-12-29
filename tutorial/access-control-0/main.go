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
	lounge := NewSpace("tcp://localhost:31417/lounge")
	go server(&lounge, &rules, &space)

	user("Alice", &lounge)
	user("Bob", &lounge)
	user("Trend", &lounge)

	space.Query("done")

}

func user(me string, lounge *Space) {
	var s string
	var t Tuple
	var err string

	lounge.Put(me, "QueryP", CreateTemplate("board", "Alice", &s, &s))
	t, _ = lounge.Get("reply", &s)
	if t.GetFieldAt(1) == "allowed" {
		t, _ = lounge.Get(me, "result", &t, &err)
		fmt.Printf("%s: I read %v on Alice's board\n", me, t.GetFieldAt(2))
	}

	lounge.Put(me, "QueryP", CreateTemplate("inbox", "Alice", &s, &s))
	t, _ = lounge.Get("reply", &s)
	if t.GetFieldAt(1) == "allowed" {
		t, _ = lounge.Get(me, "result", &t, &err)
		fmt.Printf("%s: I read %v on Alice's inbox\n", me, t.GetFieldAt(2))
	}
}

func server(lounge *Space, rules *Space, space *Space) {
	var subject string
	var action string
	var template Template
	var decision string
	for {
		t, _ := lounge.Get(&subject, &action, &template)
		subject = (t.GetFieldAt(0)).(string)
		action = (t.GetFieldAt(1)).(string)
		template = (t.GetFieldAt(2)).(Template)
		fmt.Printf("%s wants to do %s(%v)\n", subject, action, template)
		fmt.Printf("Checking rules...\n")
		decision = check(rules, subject, action, template)
		if decision != "permit" {
			fmt.Printf("Permission denied.\n")
			lounge.Put("reply", "denied")
			continue
		}
		fmt.Printf("Permission granted.\n")
		lounge.Put("reply", "allowed")
		fmt.Printf("Performing action...\n")
		switch action {
		case "Query":
			t, err := space.Query(template.Fields()...)
			lounge.Put(subject, "result", t, msg(err))
		case "QueryP":
			t, err := space.QueryP(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lounge.Put(subject, "result", t, msg(err))
		case "Get":
			t, err := space.Get(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lounge.Put(subject, "result", t, msg(err))
		case "GetP":
			t, err := space.Query(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lounge.Put(subject, "result", t, msg(err))
		case "Put":
			t, err := space.Query(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lounge.Put(subject, "result", t, msg(err))
		case "QueryAll":
			tl, err := space.QueryAll(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lounge.Put("result", tl, err)
		case "GetAll":
			tl, err := space.GetAll(template.Fields()...)
			fmt.Printf("Providing result...\n")
			lounge.Put("result", tl, err)
		default:
			fmt.Printf("Providing result...\n")
			lounge.Put("result", nil, nil)
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
