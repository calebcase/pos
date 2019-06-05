package cmd

import "github.com/spf13/cobra"

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "disk solver commands",
}

func init() {
	rootCmd.AddCommand(diskCmd)

	diskCmd.PersistentFlags().StringP("puzzle", "p", "", "Path to a puzzle config")

	diskCmd.PersistentFlags().StringP("image", "i", "", "Path to a image file")
	cobra.MarkFlagRequired(diskCmd.PersistentFlags(), "image")
}
