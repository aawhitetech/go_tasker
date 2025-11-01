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

func main() {
	fmt.Println("tasker running")

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Error: No Command Provided")
		PrintUsage()
		return
	}

	cmd:= args[0]

	switch cmd {
		case "add":
			if len(args) < 2 {
				fmt.Println("Error: No Task Name Provided")
				fmt.Println("Usage: tasker add <task description>")
				return
			}
			taskDescription := args[1]
			fmt.Printf("Added task: %s\n", taskDescription)

		case "list":
			fmt.Println("Listing tasks: (none yet)")

		default:
			fmt.Printf("Error: Unknown Command: %q\n", args[0])
			PrintUsage()
			return
	}
}
