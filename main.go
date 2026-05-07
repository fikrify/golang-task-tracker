package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Println("Usage: task-cli <command> [arguments]")
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		log.Println("Error:", err)
		return
	}

	cmd, args := os.Args[1], os.Args
	switch cmd {
	case "list":
		handleList(tasks, args)
	case "add":
		handleAdd(tasks, args)
	case "update":
		handleUpdate(tasks, args)
	default:
		if strings.HasPrefix(cmd, "mark-") {
			handleMark(tasks, args)
		} else {
			log.Println("Unknown command:", cmd)
		}
	}
}

func handleList(tasks []Task, args []string) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	if len(args) < 3 {
		for _, t := range tasks {
			printTask(t)
		}
		return
	}

	status, ok := parseStatus(args[2])
	if !ok {
		log.Println("Usage: task-cli list <todo|in-progress|done>")
		return
	}

	found := false
	for _, t := range tasks {
		if t.Status == status {
			printTask(t)
			found = true
		}
	}
	if !found {
		fmt.Println("No tasks found")
	}
}

func handleAdd(tasks []Task, args []string) {
	if len(args) < 3 {
		log.Println("Usage: task-cli add \"task name\"")
		return
	}

	now := time.Now()
	// ID is len+1; assumes no deletions. Safe for append-only task lists.
	task := Task{
		ID:          len(tasks) + 1,
		Description: args[2],
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tasks = append(tasks, task)

	if err := saveTasks(tasks); err != nil {
		log.Println("Error saving:", err)
		return
	}

	fmt.Println("Task added:", task.Description)
}

func handleUpdate(tasks []Task, args []string) {
	if len(args) < 4 {
		log.Println("Usage: task-cli update <task id> \"task name\"")
		return
	}

	taskID, err := strconv.Atoi(args[2])
	if err != nil {
		log.Println("Invalid task id")
		return
	}

	idx := findTaskIndex(tasks, taskID)
	if idx == -1 {
		log.Printf("Task %d not found", taskID)
		return
	}

	tasks[idx].Description = args[3]
	tasks[idx].UpdatedAt = time.Now()

	if err := saveTasks(tasks); err != nil {
		log.Println("Error saving:", err)
		return
	}

	fmt.Printf("Task %d updated: %s\n", taskID, tasks[idx].Description)
}

func handleMark(tasks []Task, args []string) {
	if len(args) < 3 {
		log.Println("Usage: task-cli mark-<todo|in-progress|done> <task id>")
		return
	}

	statusStr, _ := strings.CutPrefix(args[1], "mark-")
	status, ok := parseStatus(statusStr)
	if !ok {
		log.Println("Usage: task-cli mark-<todo|in-progress|done> <task id>")
		return
	}

	taskID, err := strconv.Atoi(args[2])
	if err != nil {
		log.Println("Invalid task id")
		return
	}

	idx := findTaskIndex(tasks, taskID)
	if idx == -1 {
		log.Printf("Task %d not found", taskID)
		return
	}

	tasks[idx].Status = status
	tasks[idx].UpdatedAt = time.Now()

	if err := saveTasks(tasks); err != nil {
		log.Println("Error saving:", err)
		return
	}

	fmt.Printf("Task %d marked %s\n", taskID, status)
}
