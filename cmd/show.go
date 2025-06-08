package cmd

import (
	"fmt"
	"github.com/kevin7254/task/store"
	"github.com/spf13/cobra"
	"strconv"
)

func NewShowCmd(store store.TaskRepository) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "show [ID]",
		Short: "Show (all) info about a specific task",
		Args:  cobra.MinimumNArgs(1),
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

			fmt.Println(task)

			return nil
		},
	}
	return cobraCmd
}
