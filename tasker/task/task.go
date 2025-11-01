package task

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func Add(tasks []Task, description string) []Task {
	newTask := Task{
			ID:          len(tasks) + 1,
			Description: description,
			Done:        false,
		}
	return append(tasks, newTask)
}

// MarkDone sets Done=true for a given ID and reports success.
func MarkDone(tasks []Task, id int) (bool, []Task) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			return true, tasks
		}
	}
	return false, tasks
}