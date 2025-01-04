package main

import (
	"fmt"
	"gocommandline/interactive/todo"
	"os"
	"strings"
)

// Hardcoding the file name
const todoFileName = ".todo.json"

func main() {
	l := &todo.List{}

	// Use the Get method to read todo items from file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the number of arguments provided
	switch {
	// For no extra arguments, print the list
	case len(os.Args) == 1:
		for _, item := range *l {
			fmt.Println(item.Task)
		}

		// Concatenate all provided arguments with a space and
		// add to the list as an item
	default:
		// Concatenate all arguments with a space
		item := strings.Join(os.Args[1:], " ")

		// Add the task
		l.Add(item)

		//Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
