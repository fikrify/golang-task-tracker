package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Status string

const (
	StatusPending    Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const fileName = "tasks.json"

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: task-cli <command> [arguments]")
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch args[1] {
	case "list":
		handleList(tasks)
	case "add":
		handleAdd(tasks, args)
	case "update":
		handleUpdate(tasks, args)
	case "mark-in-progress":
		handleMark(tasks, args, StatusInProgress)
	case "mark-done":
		handleMark(tasks, args, StatusDone)
	default:
		fmt.Println("Unknown command")
	}
}

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, err
	}

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

func handleList(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	for _, t := range tasks {
		fmt.Printf("[%d] %s (%s)\n", t.Id, t.Description, t.Status)
	}
}

func handleAdd(tasks []Task, args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: task-cli add \"task name\"")
		return
	}

	now := time.Now()

	task := Task{
		Id:          len(tasks) + 1,
		Description: args[2],
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, task)

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving:", err)
		return
	}

	fmt.Println("Task added:", task.Description)
}

func handleUpdate(tasks []Task, args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: task-cli update <task id> \"task name\"")
		return
	}

	taskId, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Invalid task id")
		return
	}
	taskDescription := args[3]

	for i := range tasks {
		if tasks[i].Id == taskId {
			tasks[i].Description = taskDescription
			tasks[i].UpdatedAt = time.Now()
			break
		}
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving:", err)
		return
	}

	fmt.Printf("Task %d update: %s\n", taskId, taskDescription)
}

func handleMark(tasks []Task, args []string, status Status) {
	if len(args) < 3 {
		fmt.Println("Usage: task-cli mark-in-progress <task id>")
		return
	}

	taskId, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Invalid task id")
		return
	}

	for i := range tasks {
		if tasks[i].Id == taskId {
			tasks[i].Status = status
			break
		}
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving:", err)
		return
	}

	fmt.Printf("Task %d marked in progress\n", taskId)
}
