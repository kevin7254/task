package cmd

import (
	"github.com/spf13/cobra"
	"sort"
	"strings"
	"task/data"
)

var (
	// List command flags
	projectFilter string
	showCompleted bool
	sortBy        string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long: `List tasks with optional filtering and sorting.

Examples:
  task list                  # List all incomplete tasks
  task list --all            # List all tasks including completed ones
  task list --project work   # List tasks in the 'work' project
  task list --sort priority  # Sort tasks by priority`,
	RunE: func(cmd *cobra.Command, args []string) error {
		allTasks := TaskStore.ListAllTasks()

		if len(allTasks) == 0 {
			cmd.Println("No tasks found.")
			return nil
		}

		// Filter tasks
		var filteredTasks []*data.Task
		for _, task := range allTasks {
			// Filter by completion status
			if !showCompleted && !task.CompletedAt.IsZero() {
				continue
			}

			// Filter by project
			if projectFilter != "" && !strings.EqualFold(task.Project, projectFilter) {
				continue
			}

			filteredTasks = append(filteredTasks, task)
		}

		if len(filteredTasks) == 0 {
			cmd.Println("No tasks match the filter criteria.")
			return nil
		}

		// Sort tasks
		switch strings.ToLower(sortBy) {
		case "priority":
			sort.Slice(filteredTasks, func(i, j int) bool {
				return filteredTasks[i].Priority > filteredTasks[j].Priority
			})
		case "due":
			sort.Slice(filteredTasks, func(i, j int) bool {
				return filteredTasks[i].DueDate.Before(filteredTasks[j].DueDate)
			})
		case "id":
			sort.Slice(filteredTasks, func(i, j int) bool {
				return filteredTasks[i].ID < filteredTasks[j].ID
			})
		default:
			// Default sort by ID
			sort.Slice(filteredTasks, func(i, j int) bool {
				return filteredTasks[i].ID < filteredTasks[j].ID
			})
		}

		// Display tasks
		cmd.Println("ID | Status | Priority | Due Date    | Project | Title")
		cmd.Println("---|--------|----------|-------------|---------|------------------")
		for _, task := range filteredTasks {
			status := "⏳"
			if !task.CompletedAt.IsZero() {
				status = "✅"
			} else if task.IsOverdue() {
				status = "⚠️"
			}

			priorityStr := "Low"
			if task.Priority == data.Medium {
				priorityStr = "Medium"
			} else if task.Priority == data.High {
				priorityStr = "High"
			}

			cmd.Printf("%2d | %s | %-8s | %s | %-7s | %s\n",
				task.ID,
				status,
				priorityStr,
				task.DueDate.Format("2006-01-02"),
				task.Project,
				task.Title,
			)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().StringVarP(&projectFilter, "project", "p", "", "Filter tasks by project")
	listCmd.Flags().BoolVarP(&showCompleted, "all", "a", false, "Show completed tasks")
	listCmd.Flags().StringVarP(&sortBy, "sort", "s", "id", "Sort tasks by: id, priority, or due")
}
