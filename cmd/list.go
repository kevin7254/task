package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"task/data"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all current present tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		if data.GlobalStore == nil {
			cmd.PrintErrln("Error: Store not initialized.")
			return
		}
		fmt.Println("Currently present tasks:")
		fmt.Println()
		allTasks := data.GlobalStore.ListAllTasks()

		for i := range allTasks {
			cmd.Println(allTasks[i])
		}

	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
