package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"task/model"
	"task/store"
	"time"
)

func NewAddCmd(store store.TaskRepository) *cobra.Command {
	var (
		description string
		project     string
		priority    int
		dueDate     string
	)

	addCmd := &cobra.Command{
		Use:   "add TASK_NAME",
		Short: "Add a new task",
		Long: `Add a new task to your task list.

Example:
  task add "Complete project report"
  task add "Force push to prod" --project work --priority 2 --due 2025-06-03`,
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := strings.Join(args, " ")
			if taskName == "" {
				return fmt.Errorf("task name cannot be empty")
			}

			var due time.Time
			if dueDate != "" {
				parsedDate, err := time.Parse("2006-01-02", dueDate)
				if err != nil {
					return fmt.Errorf("invalid date format: %w", err)
				}
				due = parsedDate
			} else {
				due = time.Now().AddDate(0, 0, 1)
			}

			if priority < 1 || priority > 3 {
				return fmt.Errorf("priority must be between 1 (Low) and 3 (High)")
			}

			newTask := model.NewTask(taskName, description, project, model.Priority(priority), due)

			if err := store.AddTask(newTask); err != nil {
				return fmt.Errorf("failed to add task: %w", err)
			}

			cmd.Printf("Successfully added task: %s (ID: %d)\n", newTask.Title, newTask.ID)
			return nil
		},
	}

	addCmd.Flags().StringVarP(&description, "description", "d", "", "Task description.")
	addCmd.Flags().StringVarP(&project, "project", "p", "work", "Project the task belongs to. For example work or private.")
	addCmd.Flags().IntVarP(&priority, "priority", "P", 1, "Task priority (1=Low, 2=Medium, 3=High)")
	addCmd.Flags().StringVar(&dueDate, "due", "", "Due date (format: YYYY-MM-DD)")
	return addCmd
}
