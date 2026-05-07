package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	cmd := args[1]
	switch {
	case cmd == "list":
		handleList(tasks, args)
	case cmd == "add":
		handleAdd(tasks, args)
	case cmd == "update":
		handleUpdate(tasks, args)
	case strings.Contains(cmd, "mark"):
		handleMark(tasks, args)
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

func handleList(tasks []Task, args []string) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	// Show all tasks if no status is specified
	if len(args) < 3 {
		for _, t := range tasks {
			fmt.Printf("[%d] %s (%s)\n", t.Id, t.Description, t.Status)
		}
		return
	}

	// Map status strings to Status enum values
	statusMap := map[string]Status{
		string(StatusPending):    StatusPending,
		string(StatusInProgress): StatusInProgress,
		string(StatusDone):       StatusDone,
	}

	// Check if the status is valid
	status, ok := statusMap[args[2]]
	if !ok {
		fmt.Println("Usage: task-cli list <todo|in-progress|done>")
		return
	}

	// Show tasks with the specified status
	var filtered []Task
	for _, t := range tasks {
		if t.Status == status {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) == 0 {
		fmt.Println("No tasks found")
		return
	}

	for _, t := range filtered {
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

func handleMark(tasks []Task, args []string) {
	// Map status strings to Status enum values
	statusMap := map[string]Status{
		string(StatusPending):    StatusPending,
		string(StatusInProgress): StatusInProgress,
		string(StatusDone):       StatusDone,
	}

	// Check if the status is valid
	status, ok := statusMap[strings.Replace(args[1], "mark-", "", 1)]
	if !ok || len(args) < 3 {
		fmt.Println("Usage: task-cli list <todo|in-progress|done> <task id>")
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
			tasks[i].UpdatedAt = time.Now()
			break
		}
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving:", err)
		return
	}

	fmt.Printf("Task %d marked %s \n", taskId, status)
}
