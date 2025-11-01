package main

import (
	"fmt"
	"os"
)

func PrintUsage() {
	fmt.Println("Usage: tasker <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  add <task description>")
	fmt.Println("  list")
}

type Task struct {
	ID          int
	Description string
	Done        bool
}

var tasks []Task

func main() {
	fmt.Println("tasker running")

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Error: No Command Provided")
		PrintUsage()
		return
	}

	cmd := args[0]

	switch cmd {
	case "add":
		if len(args) < 2 {
			fmt.Println("Error: No Task Name Provided")
			fmt.Println("Usage: tasker add <task description>")
			return
		}
		taskDescription := args[1]
		newTask := Task {
			ID: len(tasks) + 1,
			Description: taskDescription,
			Done: false,
		}
		tasks = append(tasks, newTask)
		fmt.Printf("Added task: %d. [%v] %s\n", newTask.ID, newTask.Done, newTask.Description)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("Listing tasks: (none yet)")
			return
		}
		for _, v := range tasks {
			fmt.Printf("%d. [%v] %s\n", v.ID, v.Done, v.Description)
		}

	default:
		fmt.Printf("Error: Unknown Command: %q\n", args[0])
		PrintUsage()
		return
	}
}
