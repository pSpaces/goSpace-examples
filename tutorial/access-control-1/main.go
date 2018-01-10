// This is a simple access control system that implements
// a simple accss control language
// see https://gitlab.gbar.dtu.dk/pSpaces/brainstorm/blob/master/secure-spaces.md

package main

import (
	"fmt"

	. "github.com/pspaces/gospace"
)

func main() {

	// Space with the data to be protected
	space := NewSpace("tcp://localhost:31415/space")
	space.Put("board", "Alice", "Charlie", "Hello")
	space.Put("inbox", "Alice", "Dave", "Some secret")

	// Space with access control rules
	rules := NewSpace("tcp://localhost:31416/policy")
	var s string
	rules.Put("policy", ">", "policy0", "policy1")

	rules.Put("policy0", "+", "policy00", "policy01")

	rules.Put("policy00", "+", "policy000", "policy001")

	rules.Put("policy000", "clause", "policy0000", "")
	rules.Put("policy0000", "permit", "Alice", "QueryP", CreateTemplate("inbox", "Alice", &s, &s))

	rules.Put("policy001", "clause", "policy0010", "")
	rules.Put("policy0010", "permit", "Alice", "QueryP", CreateTemplate("board", "Alice", &s, &s))

	rules.Put("policy01", "clause", "policy010", "")
	rules.Put("policy010", "permit", "Bob", "QueryP", CreateTemplate("board", "Alice", &s, &s))

	rules.Put("policy1", "decision", "policy11", "")
	rules.Put("policy11", "deny")

	// Space to interact with users
	lobby := NewSpace("tcp://localhost:31417/lobby")
	go server(&lobby, &rules, &space)

	user("Alice", &lobby)
	user("Bob", &lobby)
	user("Charlie", &lobby)

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
		decision = check(rules, "policy", subject, action, template)
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

func check(rules *Space, policy string, subject string, action string, template Template) string {
	var decision string
	decision = "deny"
	var s string

	t, err := rules.QueryP(policy, &s, &s, &s)
	if err != nil {
		return "maybe"
	}

	operator := t.GetFieldAt(1)
	switch operator {

	case "decision":
		policy = (t.GetFieldAt(2)).(string)
		t, err = rules.QueryP(policy, &s)
		if err != nil {
			decision = "maybe"
		} else {
			decision = (t.GetFieldAt(1)).(string)
		}
		return decision

	case "clause":
		policy = (t.GetFieldAt(2)).(string)
		t, err = rules.QueryP(policy, &s, subject, action, template)
		if err != nil {
			decision = "maybe"
		} else {
			decision = (t.GetFieldAt(1)).(string)
		}
		return decision

	case "+":
		policy = (t.GetFieldAt(2)).(string)
		decision1 := check(rules, policy, subject, action, template)
		policy = (t.GetFieldAt(3)).(string)
		decision2 := check(rules, policy, subject, action, template)
		if decision1 == decision2 {
			return decision1
		}
		if decision1 == "conflict" || decision2 == "conflict" {
			return "conflict"
		}
		if decision1 == "maybe" {
			return decision2
		}
		if decision2 == "maybe" {
			return decision1
		}
		if decision1 != decision2 {
			return "conflict"
		}
		return decision1

	case ">":
		policy = (t.GetFieldAt(2)).(string)
		decision = check(rules, policy, subject, action, template)
		if decision != "permit" && decision != "deny" {
			policy = (t.GetFieldAt(3)).(string)
			return check(rules, policy, subject, action, template)
		}

	case ">>":
		policy = (t.GetFieldAt(2)).(string)
		decision = check(rules, policy, subject, action, template)
		if decision == "maybe" {
			policy = (t.GetFieldAt(3)).(string)
			return check(rules, policy, subject, action, template)
		}

	default:
		return "maybe"

	}

	return decision
}
