package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var puzzleMaskCmd = &cobra.Command{
	Use:   "mask",
	Short: "create a mask for a puzzle",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mask called")
	},
}

func init() {
	puzzleCmd.AddCommand(puzzleMaskCmd)
}
