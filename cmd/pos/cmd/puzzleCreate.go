package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var puzzleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new puzzle",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}

func init() {
	puzzleCmd.AddCommand(puzzleCreateCmd)
}
