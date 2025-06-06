package cmd

import (
	"github.com/spf13/cobra"
	"strings"
	"task/data"
	"time"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task",
	Run: func(cmd *cobra.Command, args []string) {
		if data.GlobalStore == nil {
			cmd.PrintErrln("Error: Store not initialized.")
			return
		}

		taskName := strings.Join(args, " ")
		if taskName == "" {
			cmd.PrintErrln("Error: Task name cannot be empty.")
			return
		}

		newTask := data.NewTask(taskName, "", "pending", 0, time.Now())

		err := data.GlobalStore.AddTask(newTask)
		if err != nil {
			cmd.PrintErrf("Error adding task: %v\n", err)
			return
		}
		cmd.Printf("Successfully added task: %s\n", newTask.Title)
	},
}

func init() {
	RootCmd.AddCommand(AddCmd)
}
