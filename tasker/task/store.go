package task

import (
	"encoding/json"
	"io"
	"os"
)

const filename = "tasks.json"

func Load() ([]Task, error) {
	f, err := os.Open(filename)
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

func Save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
