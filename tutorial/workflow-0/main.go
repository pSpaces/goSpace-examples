// This is an implementation of the following simple workflow:
//
//      /---> Bob -------\
//     |								 |
//     |                 v
// Alice                Dave
//     |                 ^
//     |                 |
//      \---> Charlie ---/

package main

import (
	"fmt"
	"time"

	. "github.com/pspaces/gospace"
)

func main() {

	//runtime.GOMAXPROCS(100)

	rules := NewSpace("tcp://localhost:31415/rules")
	tokens := NewSpace("tcp://localhost:31146/tokens")

	rules.Put("<", "Alice", "Bob")
	rules.Put("<", "Alice", "Charlie")
	rules.Put("<", "Bob", "Dave")
	rules.Put("<", "Charlie", "Dave")
	rules.Put("<", "Dave", "Master")

	go coordinateTask(task, "Alice", &rules, &tokens)
	go coordinateTask(task, "Bob", &rules, &tokens)
	go coordinateTask(task, "Charlie", &rules, &tokens)
	go coordinateTask(task, "Dave", &rules, &tokens)

	tokens.Get("token", "Dave", "Master")

}

func task(me string) {
	time.Sleep(0 * time.Second)
	fmt.Printf("%s is running...\n", me)
}

func coordinateTask(task func(string), me string, rules *Space, tokens *Space) {

	// Read order constraints
	var who string
	input, _ := rules.QueryAll("<", &who, me)
	output, _ := rules.QueryAll("<", me, &who)

	// Wait for tokens of previous tasks
	for _, edge := range input {
		who = (edge.GetFieldAt(1)).(string)
		fmt.Printf("%s is waiting for %s...\n", me, who)
		tokens.Get("token", who, me)
		fmt.Printf("%s got token from %s...\n", me, who)
		//tokens, _ := space.QueryAll("token", &s1, &s2)
		//fmt.Println(tokens)
	}

	// Execute task
	task(me)

	// Send tokens to next taks
	for _, edge := range output {
		who = (edge.GetFieldAt(2)).(string)
		tokens.Put("token", me, who)
		fmt.Printf("%s sent token to %s...\n", me, who)
		//tokens, _ := space.QueryAll("token", &s1, &s2)
		//fmt.Println(tokens)
	}
}
