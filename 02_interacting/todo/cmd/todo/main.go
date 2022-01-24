package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/adamwoolhether/cliApps/02_interacting/todo"
)

const todoFileName = ".todo.json"

func main() {
	l := &todo.List{}

	// Use the Get method to read todo items from file.
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on number of args provided.
	switch {
	// Print the list for no extra args.
	case len(os.Args) == 1:
		for _, item := range *l {
			fmt.Println(item.Task)
		}
	default:
		// Concatenate all provided args with a space
		// and add to the list as an item.
		item := strings.Join(os.Args[1:], " ")
		// Add the item
		l.Add(item)
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
