package cmd_test

import (
	"bytes"
	"github.com/kevin7254/task/cmd"
	"github.com/kevin7254/task/store"
	"github.com/spf13/cobra"
	"path/filepath"
	"strconv"
	"strings"
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

	s, err := store.NewJsonStore(testStorageFile)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}
	return s
}

func TestDoCmd(t *testing.T) {
	s := setupTestStorage(t)
	cobraCmd := cmd.NewRootCmd(s)

	taskName := "My New Test Task"
	output := "Successfully added task: "
	assertCommand(t, cobraCmd, "add", taskName, output)
	assertListTasks(t, s, taskName)

	tasks := s.ListAllTasks()

	// do with ID
	doOutput := "Completed task "
	assertCommand(t, cobraCmd, "do", strconv.Itoa(tasks[0].ID), doOutput)
}

func TestAddCmd(t *testing.T) {
	s := setupTestStorage(t)
	cobraCmd := cmd.NewRootCmd(s)
	taskName := "My New Test Task"
	output := "Successfully added task: "

	assertCommand(t, cobraCmd, "add", taskName, output)
	assertListTasks(t, s, taskName)
}

func TestAddCmd_EmptyTaskName(t *testing.T) {
	s := setupTestStorage(t)
	cobraCmd := cmd.NewRootCmd(s)

	args := []string{"add"} // No task name provided
	output, err := executeCommand(cobraCmd, args...)

	if err == nil {
		t.Error("Expected error for empty task name, but got nil")
	}

	if !strings.Contains(output, "task name cannot be empty") {
		t.Errorf("Expected error message about empty task name, but got %q", output)
	}
}

func TestAddCmd_WithFlags(t *testing.T) {
	s := setupTestStorage(t)
	cobraCmd := cmd.NewRootCmd(s)

	args := []string{
		"add",
		"Test Task With Flags",
		"--description", "This is a test description",
		"--project", "test-project",
		"--priority", "2",
		"--due", "2023-12-31",
	}

	output, execErr := executeCommand(cobraCmd, args...)
	if execErr != nil {
		t.Fatalf("executeCommand failed: %v. Output: %s", execErr, output)
	}

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

func assertCommand(t testing.TB, cobraCmd *cobra.Command, inputCmd, taskName, output string) {
	t.Helper()
	outputExec, execErr := executeCommand(cobraCmd, []string{inputCmd, taskName}...)

	assertErr(t, output, execErr)
	want := output + taskName
	assertOutputContains(t, want, outputExec)
}

func assertOutputContains(t testing.TB, want, got string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Errorf("Expected output to contain task name, but got %q", got)
	}
}

func assertErr(t testing.TB, output string, execErr error) {
	t.Helper()
	if execErr != nil {
		t.Fatalf("executeCommand failed: %v. Output: %s", execErr, output)
	}
}

func assertListTasks(t testing.TB, store *store.JsonStore, taskName string) {
	t.Helper()
	tasks := store.ListAllTasks()
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
