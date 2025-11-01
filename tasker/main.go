package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
)

func PrintUsage() {
	fmt.Println("Usage: tasker <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  add <task description>")
	fmt.Println("  list")
	fmt.Println("  done <task_id>")
}

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func loadTasks() ([]Task, error) {
	f, err := os.Open("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "	")
	if err != nil {
		return err
	}
	err = os.WriteFile("tasks.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Error: No Command Provided")
		PrintUsage()
		return
	}

	cmd := args[0]

	tasks, err := loadTasks()
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
		newTask := Task{
			ID:          len(tasks) + 1,
			Description: taskDescription,
			Done:        false,
		}
		tasks = append(tasks, newTask)
		err := saveTasks(tasks)
		if err != nil {
			fmt.Printf("Error: Unable to save tasks: %s\n", err)
			return
		}
		fmt.Printf("Added task: %d. [%v] %s\n", newTask.ID, newTask.Done, newTask.Description)

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

		idx := slices.IndexFunc(tasks, func(t Task) bool { return t.ID == taskId })

		if idx == -1 {
			fmt.Printf("Error: Unable to find Task with Task Id: %d\n", taskId)
			return
		}

		tasks[idx].Done = true
		err = saveTasks(tasks)
		if err != nil {
			fmt.Printf("Error: Unable to save tasks: %s\n", err)
			return
		}
		fmt.Printf("Task marked done: %d. [%v] %s\n", tasks[idx].ID, tasks[idx].Done, tasks[idx].Description)

	default:
		fmt.Printf("Error: Unknown Command: %q\n", args[0])
		PrintUsage()
		return
	}
}
