package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adamwoolhether/cliApps/02_interacting/todo"
)

// Default file name.
var todoFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed for study.\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2022\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	// Parsing command-line flags.
	task := flag.String("task", "", "Task to be included in the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	// Check for user-defined ENV VAR to specify custom file name.
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}

	// Use the Get method to read todo items from file.
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on number of args provided.
	switch {
	case *list:
		// List current to do items
		fmt.Print(l)
	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *task != "":
		// Add the task
		l.Add(*task)
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}
