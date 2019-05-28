package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "disk solver commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("disk called")
	},
}

func init() {
	rootCmd.AddCommand(diskCmd)
}
