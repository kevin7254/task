package main

import "task/cmd"

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		return
	}
}
