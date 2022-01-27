package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adamwoolhether/cliApps/02_interacting/todo"
)

// Default file name.
var todoFileName = ".todo.json"

func main() {
	// Parsing command-line flags.
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("del", 0, "Delete a task from the todo list")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool.\nSubmit todo items at command line with the -add flag.\n"+
				"Press Enter immediately after -add to submit multiple commands.\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2022\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Check for user-defined ENV VAR to specify custom file name.
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}

	// Use the Get method to read to do items from file.
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
	case *add:
		// Any args (excluding flags) will be used as the new task.
		err := getTesk(os.Stdin, l, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err = l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *del > 0:
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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

// getTask decides where to get the description for a new task
// from: arguments of STDIN.
func getTesk(r io.Reader, l *todo.List, args ...string) error {
	if len(args) > 0 {
		l.Add(strings.Join(args, " "))
		return nil
	}

	s := bufio.NewScanner(r)
	counter := 0

	for s.Scan() {
		if err := s.Err(); err != nil {
			return err
		}

		switch {
		case counter == 0 && len(s.Text()) == 0:
			return fmt.Errorf("task cannot be blank")
		case counter > 0 && len(s.Text()) == 0:
			return nil
		}

		l.Add(s.Text())
		counter++
	}

	return nil
}
