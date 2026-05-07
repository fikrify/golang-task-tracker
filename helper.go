package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type Status string

const (
	StatusPending    Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const fileName = "tasks.json"

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, err
	}

	// Empty file is valid — treat as no tasks yet.
	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func parseStatus(s string) (Status, bool) {
	switch Status(s) {
	case StatusPending, StatusInProgress, StatusDone:
		return Status(s), true
	}
	return "", false
}

func findTaskIndex(tasks []Task, id int) int {
	for i := range tasks {
		if tasks[i].ID == id {
			return i
		}
	}
	return -1
}

func printTask(w io.Writer, t Task) error {
	_, err := fmt.Fprintln(w, fmt.Sprintf("[%d] %s (%s)", t.ID, t.Description, t.Status))
	return err
}
