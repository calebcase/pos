package cmd

import "github.com/spf13/cobra"

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream commands",
}

func init() {
	rootCmd.AddCommand(streamCmd)

	streamCmd.PersistentFlags().StringP("puzzle", "p", "", "Path to a puzzle config")
}
