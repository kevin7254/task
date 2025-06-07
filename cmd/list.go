package cmd

import (
	"fmt"
	"github.com/kevin7254/task/model"
	"github.com/kevin7254/task/store"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
	"strings"
)

var (
	projectFilter string
	showCompleted bool
	sortBy        string
	showAllInfo   string
)

type columnData struct {
	content string
	width   int
}

type taskDisplayRow struct {
	id       columnData
	status   columnData
	priority columnData
	dueDate  columnData
	project  columnData
	title    columnData
}

// buildTaskDisplayRows converts tasks into display rows and calculates optimal column widths
func buildTaskDisplayRows(tasks []*model.Task) ([]taskDisplayRow, taskDisplayRow) {
	displayRows := make([]taskDisplayRow, len(tasks))
	maxWidths := taskDisplayRow{
		id:       columnData{width: len("ID")},
		status:   columnData{width: len("Status")},
		priority: columnData{width: len("Priority")},
		dueDate:  columnData{width: len("Due Date")},
		project:  columnData{width: len("Project")},
		title:    columnData{width: len("Title")},
	}

	for i, task := range tasks {
		row := taskDisplayRow{}

		row.id.content = strconv.Itoa(task.ID)
		if len(row.id.content) > maxWidths.id.width {
			maxWidths.id.width = len(row.id.content)
		}

		row.title.content = task.Title
		if len(row.title.content) > maxWidths.title.width {
			maxWidths.title.width = len(row.title.content)
		}

		row.status.content = "⏳" // Default icon
		if !task.CompletedAt.IsZero() {
			row.status.content = "✅"
		} else if task.IsOverdue() {
			row.status.content = "⚠️"
		}
		if len(row.status.content) > maxWidths.status.width {
			maxWidths.status.width = len(row.status.content)
		}

		row.priority.content = "Low"
		if task.Priority == model.Medium {
			row.priority.content = "Medium"
		} else if task.Priority == model.High {
			row.priority.content = "High"
		}
		if len(row.priority.content) > maxWidths.priority.width {
			maxWidths.priority.width = len(row.priority.content)
		}

		row.dueDate.content = task.DueDate.Format("2006-01-02")
		if len(row.dueDate.content) > maxWidths.dueDate.width {
			maxWidths.dueDate.width = len(row.dueDate.content)
		}

		row.project.content = task.Project
		if len(row.project.content) > maxWidths.project.width {
			maxWidths.project.width = len(row.project.content)
		}

		displayRows[i] = row
	}
	return displayRows, maxWidths
}

// renderTaskTable formats and prints the task table with headers, separator, and rows
func renderTaskTable(cmd *cobra.Command, tasks []*model.Task, displayRows []taskDisplayRow, maxWidths taskDisplayRow) {
	// Create format strings for header and rows using the calculated widths
	headerFmt := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n",
		maxWidths.id.width, maxWidths.status.width, maxWidths.priority.width,
		maxWidths.dueDate.width, maxWidths.project.width, maxWidths.title.width)

	rowFmt := fmt.Sprintf("%%%dd | %%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n",
		maxWidths.id.width, maxWidths.status.width, maxWidths.priority.width,
		maxWidths.dueDate.width, maxWidths.project.width, maxWidths.title.width)

	// Create separator line components
	sepID := strings.Repeat("-", maxWidths.id.width)
	sepStatus := strings.Repeat("-", maxWidths.status.width)
	sepPriority := strings.Repeat("-", maxWidths.priority.width)
	sepDueDate := strings.Repeat("-", maxWidths.dueDate.width)
	sepProject := strings.Repeat("-", maxWidths.project.width)
	sepTitle := strings.Repeat("-", maxWidths.title.width)
	separator := fmt.Sprintf("%s-|-%s-|-%s-|-%s-|-%s-|-%s\n",
		sepID, sepStatus, sepPriority, sepDueDate, sepProject, sepTitle)

	// Print the table header
	cmd.Printf(headerFmt, "ID", "Status", "Priority", "Due Date", "Project", "Title")
	cmd.Print(separator)

	// Print the table rows
	for i, task := range tasks {
		rowContent := displayRows[i]
		cmd.Printf(rowFmt,
			task.ID, // Pass the integer ID for %d formatting
			rowContent.status.content,
			rowContent.priority.content,
			rowContent.dueDate.content,
			rowContent.project.content,
			rowContent.title.content,
		)
	}
}

func filterTasks(allTasks []*model.Task) []*model.Task {
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
		fmt.Println("No tasks match the filter criteria.")
		return nil
	}
	return filteredTasks
}

func sortTasksToPriority(filteredTasks []*model.Task) []*model.Task {
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
	return filteredTasks
}

// displayTasksTable formats and prints the list of tasks in a table.
// Column widths are adjusted based on the content and headers.
func displayTasksTable(cmd *cobra.Command, tasks []*model.Task, showBasicInfo bool) {
	displayRows, maxWidths := buildTaskDisplayRows(tasks)

	renderTaskTable(cmd, tasks, displayRows, maxWidths)
}

func NewListCmd(store store.TaskRepository) *cobra.Command {
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

			filteredTasks := filterTasks(allTasks)
			sortedAndFilteredTasks := sortTasksToPriority(filteredTasks)

			displayTasksTable(cmd, sortedAndFilteredTasks, showAllInfo == "basic")
			return nil
		},
	}

	listCmd.Flags().StringVarP(&projectFilter, "project", "p", "", "Filter tasks by project")
	listCmd.Flags().BoolVarP(&showCompleted, "completed", "c", false, "Show completed tasks")
	listCmd.Flags().StringVarP(&sortBy, "sort", "s", "id", "Sort tasks by: id, priority, or due")
	listCmd.Flags().StringVarP(&showAllInfo, "info", "i", "basic", "Show all info about tasks")
	return listCmd
}
