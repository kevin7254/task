package main

import (
	"github.com/kevin7254/task/cmd"
	"github.com/kevin7254/task/store"
	"log"
	"os"
	"path/filepath"
)

func main() {
	homeDir, osErr := os.UserHomeDir()
	if osErr != nil {
		log.Fatalf("Error getting home directory: %v\n", osErr)
	}

	storageFile := filepath.Join(homeDir, ".task", "tasks.json")
	jsonStore, storeErr := store.NewJsonStore(storageFile)
	if storeErr != nil {
		log.Fatalf("Error initializing storage: %v\n", storeErr)
	}

	rootCmdInstance := cmd.NewRootCmd(jsonStore)
	if cobraErr := rootCmdInstance.Execute(); cobraErr != nil {
		log.Fatalf("Error executing command: %v\n", cobraErr)
	}
}
