package cmd

import (
	"fmt"
	"github.com/kevin7254/task/store"
	"github.com/spf13/cobra"
	"strconv"
)

func NewEditCmd(store store.TaskRepository) *cobra.Command {
	var title string
	cobraCmd := &cobra.Command{
		Use:   "edit [ID]",
		Short: "Edit task",
		Long: `Edit a task. Right now only title (--title or -t) is supported.

Examples:
  task edit 1 --title "New title"           # Edit title of task with ID 1`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("exactly one task ID must be provided")
			}

			id, atoiErr := strconv.Atoi(args[0])
			if atoiErr != nil {
				return fmt.Errorf("invalid task ID: %s", args[0])
			}

			task := store.GetTaskByID(id)
			if task == nil {
				return fmt.Errorf("task with ID %d not found", id)
			}

			task.Title = title

			if err := store.UpdateTask(task); err != nil {
				return fmt.Errorf("failed to update task %d: %w", id, err)
			}

			cmd.Printf("Updated task with ID %d to: %s\n", id, task.Title)
			return nil
		},
	}
	cobraCmd.Flags().StringVarP(&title, "title", "t", "", "Edit task title")
	return cobraCmd
}
