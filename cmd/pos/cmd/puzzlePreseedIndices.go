package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var puzzlePreseedIndicesCmd = &cobra.Command{
	Use:   "preseed-indices",
	Short: "Compute preseed indices for a puzzle.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("preseedIndices called")
	},
}

func init() {
	puzzleCmd.AddCommand(puzzlePreseedIndicesCmd)
}
