package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var streamSolveCmd = &cobra.Command{
	Use:   "solve",
	Short: "Solve a puzzle with a stream solver.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("solve called")
	},
}

func init() {
	streamCmd.AddCommand(streamSolveCmd)
}
