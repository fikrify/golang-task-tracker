package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// inTempDir changes working dir to a temp dir for the duration of the test.
// Necessary because loadTasks/saveTasks use a hardcoded fileName.
func inTempDir(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.Chdir(orig); err != nil {
			t.Error(err)
		}
	}
}

func seedTasks(t *testing.T, tasks []Task) {
	t.Helper()
	if err := saveTasks(tasks); err != nil {
		t.Fatal(err)
	}
}

func makeTasks(n int) []Task {
	now := time.Now()
	tasks := make([]Task, n)
	for i := range tasks {
		tasks[i] = Task{
			ID:          i + 1,
			Description: "Task " + string(rune('A'+i)),
			Status:      StatusPending,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}
	return tasks
}

// --- handleAdd ---

func TestHandleAdd(t *testing.T) {
	defer inTempDir(t)()

	var buf bytes.Buffer
	if err := handleAdd(&buf, []Task{}, []string{"task-cli", "add", "Buy milk"}); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "Buy milk") {
		t.Errorf("output missing task name; got %q", buf.String())
	}

	tasks, _ := loadTasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task saved, got %d", len(tasks))
	}
	if tasks[0].Description != "Buy milk" {
		t.Errorf("description = %q; want %q", tasks[0].Description, "Buy milk")
	}
	if tasks[0].Status != StatusPending {
		t.Errorf("status = %q; want %q", tasks[0].Status, StatusPending)
	}
	if tasks[0].ID != 1 {
		t.Errorf("id = %d; want 1", tasks[0].ID)
	}
}

func TestHandleAddMissingArg(t *testing.T) {
	defer inTempDir(t)()

	if err := handleAdd(io.Discard, []Task{}, []string{"task-cli", "add"}); err == nil {
		t.Error("expected error for missing arg")
	}

	tasks, _ := loadTasks()
	if len(tasks) != 0 {
		t.Errorf("expected no task saved on missing arg, got %d", len(tasks))
	}
}

// --- handleUpdate ---

func TestHandleUpdate(t *testing.T) {
	defer inTempDir(t)()

	initial := makeTasks(1)
	seedTasks(t, initial)

	var buf bytes.Buffer
	if err := handleUpdate(&buf, initial, []string{"task-cli", "update", "1", "New desc"}); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "New desc") {
		t.Errorf("output missing new description; got %q", buf.String())
	}

	tasks, _ := loadTasks()
	if tasks[0].Description != "New desc" {
		t.Errorf("description = %q; want %q", tasks[0].Description, "New desc")
	}
}

func TestHandleUpdateNotFound(t *testing.T) {
	defer inTempDir(t)()

	initial := makeTasks(1)
	seedTasks(t, initial)

	if err := handleUpdate(io.Discard, initial, []string{"task-cli", "update", "99", "Ghost"}); err == nil {
		t.Error("expected error for not found task")
	}

	tasks, _ := loadTasks()
	if tasks[0].Description != initial[0].Description {
		t.Errorf("task modified unexpectedly: %q", tasks[0].Description)
	}
}

// --- handleMark ---

func TestHandleMark(t *testing.T) {
	cases := []struct {
		cmd    string
		status Status
	}{
		{"mark-done", StatusDone},
		{"mark-in-progress", StatusInProgress},
		{"mark-todo", StatusPending},
	}

	for _, tc := range cases {
		t.Run(tc.cmd, func(t *testing.T) {
			defer inTempDir(t)()

			initial := makeTasks(1)
			seedTasks(t, initial)

			var buf bytes.Buffer
			if err := handleMark(&buf, initial, []string{"task-cli", tc.cmd, "1"}); err != nil {
				t.Fatal(err)
			}

			if !strings.Contains(buf.String(), string(tc.status)) {
				t.Errorf("output missing status %q; got %q", tc.status, buf.String())
			}

			tasks, _ := loadTasks()
			if tasks[0].Status != tc.status {
				t.Errorf("status = %q; want %q", tasks[0].Status, tc.status)
			}
		})
	}
}

func TestHandleMarkNotFound(t *testing.T) {
	defer inTempDir(t)()

	initial := makeTasks(1)
	seedTasks(t, initial)

	if err := handleMark(io.Discard, initial, []string{"task-cli", "mark-done", "99"}); err == nil {
		t.Error("expected error for not found task")
	}

	tasks, _ := loadTasks()
	if tasks[0].Status != StatusPending {
		t.Errorf("status changed unexpectedly: %q", tasks[0].Status)
	}
}

// --- handleList ---

func TestHandleListAll(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task A", Status: StatusPending},
		{ID: 2, Description: "Task B", Status: StatusDone},
	}

	var buf bytes.Buffer
	if err := handleList(&buf, tasks, []string{"task-cli", "list"}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()

	if !strings.Contains(out, "Task A") || !strings.Contains(out, "Task B") {
		t.Errorf("expected both tasks in output; got %q", out)
	}
}

func TestHandleListFiltered(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task A", Status: StatusPending},
		{ID: 2, Description: "Task B", Status: StatusDone},
	}

	var buf bytes.Buffer
	if err := handleList(&buf, tasks, []string{"task-cli", "list", "done"}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()

	if !strings.Contains(out, "Task B") {
		t.Errorf("expected Task B in filtered output; got %q", out)
	}
	if strings.Contains(out, "Task A") {
		t.Errorf("Task A should be filtered out; got %q", out)
	}
}

func TestHandleListEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := handleList(&buf, []Task{}, []string{"task-cli", "list"}); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "No tasks") {
		t.Errorf("expected 'No tasks' message; got %q", buf.String())
	}
}

func TestHandleListFilterNoMatch(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task A", Status: StatusPending},
	}

	var buf bytes.Buffer
	if err := handleList(&buf, tasks, []string{"task-cli", "list", "done"}); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "No tasks") {
		t.Errorf("expected 'No tasks' for empty filter result; got %q", buf.String())
	}
}
