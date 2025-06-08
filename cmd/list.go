package cmd

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/kevin7254/task/model"
	"github.com/kevin7254/task/store"
	"github.com/spf13/cobra"
)

// listOptions holds all the flag-related values for the list command.
type listOptions struct {
	projectFilter string
	showCompleted bool
	sortBy        string
	view          string
}

// NewListCmd creates and configures the 'list' command.
func NewListCmd(taskStore store.TaskRepository) *cobra.Command {
	opts := &listOptions{}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Long: `List tasks with optional filtering and sorting.

Examples:
  task list              # List all incomplete tasks (basic view)
  task list --view full  # List all incomplete tasks (full view)
  task list -c           # List all tasks including completed ones
  task list -p work      # List tasks in the 'work' project
  task list -s priority  # Sort tasks by priority`,
		RunE: func(cmd *cobra.Command, args []string) error {
			allTasks := taskStore.ListAllTasks()

			if len(allTasks) == 0 {
				cmd.Println("No tasks found.")
				return nil
			}

			filteredTasks := filterTasks(allTasks, opts)
			if len(filteredTasks) == 0 {
				cmd.Println("No tasks match the filter criteria.")
				return nil
			}

			sortTasks(filteredTasks, opts)

			dm := NewDisplayManager(cmd.OutOrStdout())
			return dm.RenderTasks(filteredTasks, opts.view)
		},
	}

	listCmd.Flags().StringVarP(&opts.projectFilter, "project", "p", "", "Filter tasks by project")
	listCmd.Flags().BoolVarP(&opts.showCompleted, "completed", "c", false, "Show completed tasks")
	listCmd.Flags().StringVarP(&opts.sortBy, "sort", "s", "id", "Sort tasks by: id, priority, or due")
	listCmd.Flags().StringVar(&opts.view, "view", "basic", "Set view format: basic or full")

	return listCmd
}

// filterTasks returns a new slice of tasks that match the filter criteria in opts.
func filterTasks(tasks []*model.Task, opts *listOptions) []*model.Task {
	filtered := make([]*model.Task, 0, len(tasks))
	for _, task := range tasks {
		if !opts.showCompleted && !task.CompletedAt.IsZero() {
			continue
		}

		if opts.projectFilter != "" && !strings.EqualFold(task.Project, opts.projectFilter) {
			continue
		}

		filtered = append(filtered, task)
	}
	return filtered
}

// sortTasks sorts the slice of tasks in-place based on the sortBy option.
func sortTasks(tasks []*model.Task, opts *listOptions) {
	switch strings.ToLower(opts.sortBy) {
	case "priority":
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Priority > tasks[j].Priority
		})
	case "due":
		sort.Slice(tasks, func(i, j int) bool {
			// Handle tasks without due dates by sorting them last.
			if tasks[i].DueDate.IsZero() {
				return false
			}
			if tasks[j].DueDate.IsZero() {
				return true
			}
			return tasks[i].DueDate.Before(tasks[j].DueDate)
		})
	case "id":
		fallthrough // Fallthrough to a default case
	default:
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].ID < tasks[j].ID
		})
	}
}

// DisplayManager handles the rendering of data to an output stream.
type DisplayManager struct {
	writer io.Writer
}

// NewDisplayManager creates a new display manager.
func NewDisplayManager(w io.Writer) *DisplayManager {
	return &DisplayManager{writer: w}
}

// RenderTasks orchestrates the conversion of tasks to a tabular format and prints them.
func (dm *DisplayManager) RenderTasks(tasks []*model.Task, view string) error {
	headers, rows := buildTableData(tasks, view)
	if len(rows) == 0 {
		return nil // Nothing to render
	}
	dm.renderTable(headers, rows)
	return nil
}

// buildTableData transforms tasks into headers and rows based on the selected view.
func buildTableData(tasks []*model.Task, view string) (headers []string, rows [][]string) {
	switch view {
	case "basic":
		headers = []string{"ID", "Title"}
		rows = make([][]string, len(tasks))
		for i, task := range tasks {
			rows[i] = []string{
				strconv.Itoa(task.ID),
				task.Title,
			}
		}
	default: // "full" view
		headers = []string{"ID", "Status", "Priority", "Due Date", "Project", "Title"}
		rows = make([][]string, len(tasks))
		for i, task := range tasks {
			rows[i] = []string{
				strconv.Itoa(task.ID),
				getStatusIcon(task),
				getPriorityString(task.Priority),
				task.DueDate.Format("2006-01-02"),
				task.Project,
				task.Title,
			}
		}
	}
	return headers, rows
}

// renderTable is a generic function that can print any table given headers and rows.
func (dm *DisplayManager) renderTable(headers []string, rows [][]string) {
	if len(headers) == 0 || len(rows) == 0 {
		return
	}

	// 1. Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// 2. Create a format string
	var fmtBuilder strings.Builder
	for _, w := range widths {
		fmtBuilder.WriteString(fmt.Sprintf("%%-%ds | ", w))
	}
	// Trim the last "|"
	format := strings.TrimSuffix(fmtBuilder.String(), " | ") + "\n"

	// 3. Print Header
	headerContent := make([]interface{}, len(headers))
	for i, h := range headers {
		headerContent[i] = h
	}
	_, _ = fmt.Fprintf(dm.writer, format, headerContent...)

	// 4. Print Separator
	var sepBuilder strings.Builder
	for _, w := range widths {
		sepBuilder.WriteString(strings.Repeat("-", w))
		sepBuilder.WriteString("-+-")
	}
	separator := strings.TrimSuffix(sepBuilder.String(), "-+-") + "\n"
	_, _ = fmt.Fprint(dm.writer, separator)

	// 5. Print Rows
	for _, row := range rows {
		rowContent := make([]interface{}, len(row))
		for i, v := range row {
			rowContent[i] = v
		}
		_, _ = fmt.Fprintf(dm.writer, format, rowContent...)
	}
}

func getStatusIcon(task *model.Task) string {
	if !task.CompletedAt.IsZero() {
		return "✅"
	}
	if task.IsOverdue() {
		return "⚠️"
	}
	return "⏳"
}

func getPriorityString(p model.Priority) string {
	switch p {
	case model.High:
		return "High"
	case model.Medium:
		return "Medium"
	default:
		return "Low"
	}
}
