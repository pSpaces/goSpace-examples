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

	// Tuple space for the rules regulating the order of tasks
	rules := NewSpace("tcp://localhost:31415/rules")
	rules.Put("Alice", "Bob")
	rules.Put("Alice", "Charlie")
	rules.Put("Bob", "Dave")
	rules.Put("Charlie", "Dave")
	rules.Put("Dave", "Master")

	// Tuple space for the tokens used to signal termination
	tokens := NewSpace("tcp://localhost:31146/tokens")

	// Launch all task coordinators
	// Note that the first argument is the task as a closure
	go coordinateTask(func() { task("Alice") }, "Alice", &rules, &tokens)
	go coordinateTask(func() { task("Bob") }, "Bob", &rules, &tokens)
	go coordinateTask(func() { task("Charlie") }, "Charlie", &rules, &tokens)
	go coordinateTask(func() { task("Dave") }, "Dave", &rules, &tokens)

	// Wait for the Dave to finish
	tokens.Get("Dave", "Master")

}

// All task execute the same code
func task(me string) {
	fmt.Printf("%s is running...\n", me)
	time.Sleep(1 * time.Second)
	fmt.Printf("%s done!\n", me)
}

func coordinateTask(task func(), me string, rules *Space, tokens *Space) {

	// Read order constraints
	var who string
	before, _ := rules.QueryAll(&who, me)
	after, _ := rules.QueryAll(me, &who)

	// Wait for tokens of previous tasks
	for _, edge := range before {
		who = (edge.GetFieldAt(0)).(string)
		fmt.Printf("%s is waiting for %s...\n", me, who)
		tokens.Get(who, me)
		fmt.Printf("%s got token from %s...\n", me, who)
	}

	// Execute task
	task()

	// Send tokens to tasks that come next
	for _, edge := range after {
		who = (edge.GetFieldAt(1)).(string)
		tokens.Put(me, who)
		fmt.Printf("%s sent token to %s...\n", me, who)
	}
}
