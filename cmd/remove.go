package cmd

import (
	"fmt"
	"github.com/kevin7254/task/store"
	"github.com/spf13/cobra"
	"strconv"
)

func NewRemoveCmd(store store.TaskRepository) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "remove ID [ID...]",
		Short: "Remove task(s)",
		Long: `Remove one or more tasks totally. This is different
compared to "task do" in that this removes them totally, they will not
included in any stats in any way.

Examples:
  task remove 1           # Remove task with ID 1 totally
  task remove 1 2 3       # Remove multiple totally`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var ids []int
			for _, arg := range args {
				id, err := strconv.Atoi(arg)
				if err != nil {
					return fmt.Errorf("invalid task ID: %s", arg)
				}
				ids = append(ids, id)
			}

			for _, id := range ids {
				task := store.GetTaskByID(id)
				if task == nil {
					return fmt.Errorf("task with ID %d not found", id)
				}

				if err := store.DeleteTask(id); err != nil {
					return fmt.Errorf("failed to remove task %d: %w", id, err)
				}

				cmd.Printf("Remove task %d: %s\n", id, task.Title)
			}
			return nil
		},
	}
	return cobraCmd
}
