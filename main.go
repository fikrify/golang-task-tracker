package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatal("Usage: task-cli <command> [arguments]")
	}

	tasks, err := loadTasks()
	if err != nil {
		log.Fatal("error:", err)
	}

	cmd, args := os.Args[1], os.Args
	var runErr error
	switch cmd {
	case "list":
		runErr = handleList(os.Stdout, tasks, args)
	case "add":
		runErr = handleAdd(os.Stdout, tasks, args)
	case "update":
		runErr = handleUpdate(os.Stdout, tasks, args)
	default:
		if strings.HasPrefix(cmd, "mark-") {
			runErr = handleMark(os.Stdout, tasks, args)
		} else {
			log.Fatal("unknown command:", cmd)
		}
	}
	if runErr != nil {
		log.Fatal(runErr)
	}
}

func handleList(w io.Writer, tasks []Task, args []string) error {
	if len(tasks) == 0 {
		_, err := fmt.Fprintln(w, "No tasks found")
		return err
	}

	if len(args) < 3 {
		for _, t := range tasks {
			if err := printTask(w, t); err != nil {
				return err
			}
		}
		return nil
	}

	status, ok := parseStatus(args[2])
	if !ok {
		return fmt.Errorf("usage: task-cli list <todo|in-progress|done>")
	}

	found := false
	for _, t := range tasks {
		if t.Status == status {
			if err := printTask(w, t); err != nil {
				return err
			}
			found = true
		}
	}
	if !found {
		_, err := fmt.Fprintln(w, "No tasks found")
		return err
	}
	return nil
}

func handleAdd(w io.Writer, tasks []Task, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: task-cli add \"task name\"")
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
		return fmt.Errorf("saving: %w", err)
	}

	_, err := fmt.Fprintln(w, "Task added:", task.Description)
	return err
}

func handleUpdate(w io.Writer, tasks []Task, args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: task-cli update <task id> \"task name\"")
	}

	taskID, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid task id")
	}

	idx := findTaskIndex(tasks, taskID)
	if idx == -1 {
		return fmt.Errorf("task %d not found", taskID)
	}

	tasks[idx].Description = args[3]
	tasks[idx].UpdatedAt = time.Now()

	if err := saveTasks(tasks); err != nil {
		return fmt.Errorf("saving: %w", err)
	}

	_, err = fmt.Fprintln(w, fmt.Sprintf("Task %d updated: %s", taskID, tasks[idx].Description))
	return err
}

func handleMark(w io.Writer, tasks []Task, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: task-cli mark-<todo|in-progress|done> <task id>")
	}

	statusStr, _ := strings.CutPrefix(args[1], "mark-")
	status, ok := parseStatus(statusStr)
	if !ok {
		return fmt.Errorf("usage: task-cli mark-<todo|in-progress|done> <task id>")
	}

	taskID, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid task id")
	}

	idx := findTaskIndex(tasks, taskID)
	if idx == -1 {
		return fmt.Errorf("task %d not found", taskID)
	}

	tasks[idx].Status = status
	tasks[idx].UpdatedAt = time.Now()

	if err := saveTasks(tasks); err != nil {
		return fmt.Errorf("saving: %w", err)
	}

	_, err = fmt.Fprintln(w, fmt.Sprintf("Task %d marked %s", taskID, status))
	return err
}
