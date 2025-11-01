package main

import (
	"fmt"
	"os"
	"strconv"

	"tasker/task"
)

func PrintUsage() {
	fmt.Println("Usage: tasker <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  add <task description>")
	fmt.Println("  list")
	fmt.Println("  done <task_id>")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Error: No Command Provided")
		PrintUsage()
		return
	}

	cmd := args[0]

	tasks, err := task.Load()
	if err != nil {
		fmt.Printf("Error: Unable to load tasks: %s\n", err)
		return
	}

	switch cmd {
	case "add":
		if len(args) < 2 {
			fmt.Println("Error: No Task Name Provided")
			PrintUsage()
			return
		}
		taskDescription := args[1]

		tasks = task.Add(tasks, taskDescription)

		err := task.Save(tasks)
		if err != nil {
			fmt.Printf("Error: Unable to save tasks: %s\n", err)
			return
		}
		fmt.Printf("Added task: %d. %s\n", tasks[len(tasks) - 1].ID, taskDescription)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("Listing tasks: (none yet)")
			return
		}
		for _, v := range tasks {
			fmt.Printf("%d. [%v] %s\n", v.ID, v.Done, v.Description)
		}

	case "done":
		if len(args) < 2 {
			fmt.Println("Error: No Task Id Provided")
			PrintUsage()
			return
		}
		taskId, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Error: Unable to convert Task Id to integer: %s\n", err)
			return
		}

		if taskId <= 0 {
			fmt.Printf("Error: Invalid Task Id: %d\n", taskId)
			return
		}

		result, tasks := task.MarkDone(tasks, taskId)

		if result == false {
			fmt.Printf("Error: Unable to mark Task done with Task Id: %d\n", taskId)
			return
		}

		err = task.Save(tasks)
		if err != nil {
			fmt.Printf("Error: Unable to save tasks: %s\n", err)
			return
		}
		fmt.Printf("Task marked done for Task Id: %d.\n", taskId)

	default:
		fmt.Printf("Error: Unknown Command: %q\n", args[0])
		PrintUsage()
		return
	}
}
