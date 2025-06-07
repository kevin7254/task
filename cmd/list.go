package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
	"strings"
	"task/model"
	"task/store"
)

// displayTasksTable formats and prints the list of tasks in a table.
// Column widths are adjusted based on the content and headers.
func displayTasksTable(cmd *cobra.Command, tasks []*model.Task) {
	// Note: The caller (RunE) already checks if tasks is empty
	// and prints "No tasks match the filter criteria."
	// So, this function assumes 'tasks' is not empty.

	// Initial column widths based on headers
	idWidth := len("ID")
	statusWidth := len("Status")     // Header "Status" (6 chars) vs icons like "⏳" (3 bytes)
	priorityWidth := len("Priority") // Header "Priority" (8 chars) vs "Medium" (6 chars)
	dueDateWidth := len("Due Date")  // Header "Due Date" (8 chars) vs "2006-01-02" (10 chars)
	projectWidth := len("Project")
	titleWidth := len("Title")

	// Pre-calculate display strings and determine max widths
	type taskDisplayRow struct {
		idString   string
		statusIcon string
		priority   string
		dueDate    string
		project    string
		title      string
	}
	displayRows := make([]taskDisplayRow, len(tasks))

	for i, task := range tasks {
		row := taskDisplayRow{}

		row.idString = strconv.Itoa(task.ID)
		if len(row.idString) > idWidth {
			idWidth = len(row.idString)
		}

		row.statusIcon = "⏳" // Default icon
		if !task.CompletedAt.IsZero() {
			row.statusIcon = "✅"
		} else if task.IsOverdue() {
			row.statusIcon = "⚠️"
		}
		// Status icons (e.g., "⏳") are 3 bytes. "Status" header is 6 characters.
		// The header width will likely dominate unless an icon string is unexpectedly long.
		if len(row.statusIcon) > statusWidth {
			statusWidth = len(row.statusIcon)
		}

		row.priority = "Low"
		if task.Priority == model.Medium {
			row.priority = "Medium"
		} else if task.Priority == model.High {
			row.priority = "High"
		}
		if len(row.priority) > priorityWidth {
			priorityWidth = len(row.priority)
		}

		row.dueDate = task.DueDate.Format("2006-01-02")
		if len(row.dueDate) > dueDateWidth { // Ensures width accommodates "YYYY-MM-DD"
			dueDateWidth = len(row.dueDate)
		}

		row.project = task.Project
		if len(row.project) > projectWidth {
			projectWidth = len(row.project)
		}

		row.title = task.Title
		if len(row.title) > titleWidth {
			titleWidth = len(row.title)
		}
		displayRows[i] = row
	}

	// Create format strings for header and rows
	// Headers are left-aligned strings.
	headerFmt := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n",
		idWidth, statusWidth, priorityWidth, dueDateWidth, projectWidth, titleWidth)

	// Rows: ID is a right-aligned integer. Other fields are left-aligned strings.
	rowFmt := fmt.Sprintf("%%%dd | %%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n",
		idWidth, statusWidth, priorityWidth, dueDateWidth, projectWidth, titleWidth)

	// Create separator line components
	sepID := strings.Repeat("-", idWidth)
	sepStatus := strings.Repeat("-", statusWidth)
	sepPriority := strings.Repeat("-", priorityWidth)
	sepDueDate := strings.Repeat("-", dueDateWidth)
	sepProject := strings.Repeat("-", projectWidth)
	sepTitle := strings.Repeat("-", titleWidth)
	separator := fmt.Sprintf("%s-|-%s-|-%s-|-%s-|-%s-|-%s\n",
		sepID, sepStatus, sepPriority, sepDueDate, sepProject, sepTitle)

	// Print the table header
	cmd.Printf(headerFmt, "ID", "Status", "Priority", "Due Date", "Project", "Title")
	cmd.Print(separator)

	// Print the table rows
	for i, task := range tasks { // Iterate original tasks for task.ID (int)
		rowContent := displayRows[i]
		cmd.Printf(rowFmt,
			task.ID, // Pass the integer ID for %d formatting
			rowContent.statusIcon,
			rowContent.priority,
			rowContent.dueDate,
			rowContent.project,
			rowContent.title,
		)
	}
}

func NewListCmd(store store.TaskRepository) *cobra.Command {
	var (
		projectFilter string
		showCompleted bool
		sortBy        string
	)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Long: `List tasks with optional filtering and sorting.

Examples:
  task list                  # List all incomplete tasks
  task list --all            # List all tasks including completed ones
  task list --project work   # List tasks in the 'work' project
  task list --sort priority  # Sort tasks by priority`,
		RunE: func(cmd *cobra.Command, args []string) error {
			allTasks := store.ListAllTasks()

			if len(allTasks) == 0 {
				cmd.Println("No tasks found.")
				return nil
			}

			var filteredTasks []*model.Task
			for _, task := range allTasks {
				if !showCompleted && !task.CompletedAt.IsZero() {
					continue
				}

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

			// Display tasks using the new helper function
			displayTasksTable(cmd, filteredTasks)
			return nil
		},
	}

	listCmd.Flags().StringVarP(&projectFilter, "project", "p", "", "Filter tasks by project")
	listCmd.Flags().BoolVarP(&showCompleted, "all", "a", false, "Show completed tasks")
	listCmd.Flags().StringVarP(&sortBy, "sort", "s", "id", "Sort tasks by: id, priority, or due")
	return listCmd
}
