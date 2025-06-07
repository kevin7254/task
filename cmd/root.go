package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"task/data"
	"task/storage"
)

// TaskStore is the store used by all commands
var TaskStore data.StoreInterface

var (
	// RootCmd is the root command for the task CLI
	RootCmd = &cobra.Command{
		Use:   "task",
		Short: "Task is a CLI tool for managing tasks",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip initialization for help command
			if cmd.Name() == "help" {
				return nil
			}

			// Initialize storage if not already done
			if TaskStore == nil {
				store, err := initStorage()
				if err != nil {
					return err
				}
				TaskStore = store
			}
			return nil
		},
	}
)

// initStorage initializes the task storage
func initStorage() (data.StoreInterface, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	storageFile := filepath.Join(homeDir, ".task", "tasks.json")
	store, err := storage.NewStorage(storageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	return store, nil
}
