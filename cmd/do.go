package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"task/store"
)

func NewDoCmd(store store.TaskRepository) *cobra.Command {
	var timeSpent int
	cobraCmd := &cobra.Command{
		Use:   "do ID [ID...]",
		Short: "Mark task(s) as completed",
		Long: `Mark one or more tasks as completed by their IDs.

Examples:
  task do 1           # Mark task with ID 1 as completed
  task do 1 2 3       # Mark multiple tasks as completed
  task do 1 --time 30 # Mark task as completed and log 30 minutes spent`,
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

				if timeSpent > 0 {
					task.AddTimeSpent(int64(timeSpent))
				}

				task.Complete()

				if err := store.UpdateTask(task); err != nil {
					return fmt.Errorf("failed to update task %d: %w", id, err)
				}

				cmd.Printf("Completed task %d: %s\n", id, task.Title)
			}
			return nil
		},
	}
	cobraCmd.Flags().IntVarP(&timeSpent, "time", "t", 0, "Time spent on the task in minutes")
	return cobraCmd
}
