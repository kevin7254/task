package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"task/data"
	"task/storage"
)

var (
	RootCmd = &cobra.Command{
		Use:   "task",
		Short: "Task is a CLI tool for managing tasks",
	}
)

func init() {
	cobra.OnInitialize(initStorage)
}

func initStorage() {
	if data.GlobalStore != nil {
		return
	}
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("wrong dirr for: %v", err)
		os.Exit(1)
	}
	storageFile := filepath.Join(currentDir, ".task", "tasks.json")
	store, err := storage.NewStorage(storageFile)
	if err != nil {
		fmt.Printf("wrong with creating new store %v", err)
		os.Exit(1)
	}
	data.SetStore(store)
}
