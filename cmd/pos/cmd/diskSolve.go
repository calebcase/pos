package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var diskSolveCmd = &cobra.Command{
	Use:   "solve",
	Short: "Solve a puzzle with a disk solver.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("solve called")
	},
}

func init() {
	diskCmd.AddCommand(diskSolveCmd)
}
