package cmd

import (
	"github.com/kevin7254/task/store"
	"github.com/spf13/cobra"
)

func NewRootCmd(store store.TaskRepository) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "task",
		Short: "Task is a CLI tool for managing tasks",
	}
	rootCmd.AddCommand(NewAddCmd(store))
	rootCmd.AddCommand(NewDoCmd(store))
	rootCmd.AddCommand(NewListCmd(store))
	rootCmd.AddCommand(NewRemoveCmd(store))
	rootCmd.AddCommand(NewEditCmd(store))
	rootCmd.AddCommand(NewShowCmd(store))
	return rootCmd
}
