package cmd

import "github.com/spf13/cobra"

var puzzleCmd = &cobra.Command{
	Use:   "puzzle",
	Short: "puzzle commands",
}

func init() {
	rootCmd.AddCommand(puzzleCmd)
}
