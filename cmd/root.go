package cmd

import (
	"github.com/spf13/cobra"
	"task/store"
)

func NewRootCmd(store store.TaskRepository) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "task",
		Short: "Task is a CLI tool for managing tasks",
	}
	rootCmd.AddCommand(NewAddCmd(store))
	rootCmd.AddCommand(NewDoCmd(store))
	rootCmd.AddCommand(NewListCmd(store))
	return rootCmd
}
