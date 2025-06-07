package cmd_test

import (
	"bytes"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
	"task/cmd"
	"task/store"
	"testing"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	_, err = root.ExecuteC()
	return buf.String(), err
}

func setupTestStorage(t *testing.T) *store.JsonStore {
	tempDir := t.TempDir()
	testStorageFile := filepath.Join(tempDir, "test_tasks.json")

	// Initialize a real store instance pointing to the temporary file
	s, err := store.NewJsonStore(testStorageFile)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}
	return s
}

func TestAddCmd_Success(t *testing.T) {
	// 1. Setup test environment
	s := setupTestStorage(t)
	cobraCmd := cmd.NewRootCmd(s)

	// 2. Define test case
	taskName := "My New Test Task"
	args := []string{"add", taskName} // Command is "add", argument is the task name

	// 3. Execute the command
	output, execErr := executeCommand(cobraCmd, args...)
	if execErr != nil {
		t.Fatalf("executeCommand failed: %v. Output: %s", execErr, output)
	}

	// 4. Verify command output
	if !strings.Contains(output, "Successfully added task: "+taskName) {
		t.Errorf("Expected output to contain task name, but got %q", output)
	}

	// 5. Verify store state
	tasks := s.ListAllTasks()
	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task in store, found %d", len(tasks))
	}
	if tasks[0].Title != taskName {
		t.Errorf("Expected task name to be %q, got %q", taskName, tasks[0].Title)
	}
	if tasks[0].ID <= 0 {
		t.Errorf("Expected task ID to be positive, got %d", tasks[0].ID)
	}
}

func TestAddCmd_EmptyTaskName(t *testing.T) {
	// 1. Setup test environment
	s := setupTestStorage(t)
	cobraCmd := cmd.NewRootCmd(s)

	// 2. Execute command with no task name
	args := []string{"add"} // No task name provided
	output, err := executeCommand(cobraCmd, args...)

	// 3. Verify error is returned
	if err == nil {
		t.Error("Expected error for empty task name, but got nil")
	}

	// 4. Verify error message
	if !strings.Contains(output, "task name cannot be empty") {
		t.Errorf("Expected error message about empty task name, but got %q", output)
	}
}

// TestAddCmd_WithFlags tests adding a task with various flags
func TestAddCmd_WithFlags(t *testing.T) {
	// 1. Setup test environment
	s := setupTestStorage(t)
	cobraCmd := cmd.NewRootCmd(s)

	// 2. Define test case with flags
	args := []string{
		"add",
		"Test Task With Flags",
		"--description", "This is a test description",
		"--project", "test-project",
		"--priority", "2",
		"--due", "2023-12-31",
	}

	// 3. Execute the command
	output, execErr := executeCommand(cobraCmd, args...)
	if execErr != nil {
		t.Fatalf("executeCommand failed: %v. Output: %s", execErr, output)
	}

	// 4. Verify store state
	tasks := s.ListAllTasks()
	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task in store, found %d", len(tasks))
	}

	task := tasks[0]
	if task.Title != "Test Task With Flags" {
		t.Errorf("Expected task title to be 'Test Task With Flags', got %q", task.Title)
	}
	if task.Description != "This is a test description" {
		t.Errorf("Expected task description to be set, got %q", task.Description)
	}
	if task.Project != "test-project" {
		t.Errorf("Expected task project to be 'test-project', got %q", task.Project)
	}
	if task.Priority != 2 {
		t.Errorf("Expected task priority to be 2, got %d", task.Priority)
	}
}
