package main

import (
	"os"
	"testing"
	"time"
)

func TestParseStatus(t *testing.T) {
	tests := []struct {
		input string
		want  Status
		ok    bool
	}{
		{"todo", StatusPending, true},
		{"in-progress", StatusInProgress, true},
		{"done", StatusDone, true},
		{"invalid", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		got, ok := parseStatus(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("parseStatus(%q) = %q, %v; want %q, %v", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestFindTaskIndex(t *testing.T) {
	tasks := []Task{{ID: 1}, {ID: 3}, {ID: 5}}

	if i := findTaskIndex(tasks, 3); i != 1 {
		t.Errorf("findTaskIndex(3) = %d; want 1", i)
	}
	if i := findTaskIndex(tasks, 99); i != -1 {
		t.Errorf("findTaskIndex(99) = %d; want -1", i)
	}
	if i := findTaskIndex([]Task{}, 1); i != -1 {
		t.Errorf("findTaskIndex on empty = %d; want -1", i)
	}
}

func TestLoadTasksNoFile(t *testing.T) {
	defer inTempDir(t)()

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty slice, got %d tasks", len(tasks))
	}
}

func TestLoadTasksEmptyFile(t *testing.T) {
	defer inTempDir(t)()

	if err := os.WriteFile(fileName, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty slice for empty file, got %d tasks", len(tasks))
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	defer inTempDir(t)()

	now := time.Now().Truncate(time.Second)
	want := []Task{
		{ID: 1, Description: "Buy milk", Status: StatusPending, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Description: "Walk dog", Status: StatusDone, CreatedAt: now, UpdatedAt: now},
	}

	if err := saveTasks(want); err != nil {
		t.Fatal(err)
	}

	got, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].ID != want[i].ID || got[i].Description != want[i].Description || got[i].Status != want[i].Status {
			t.Errorf("task[%d] mismatch: got %+v, want %+v", i, got[i], want[i])
		}
	}
}
