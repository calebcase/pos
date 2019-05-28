package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var diskPrepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Prepare a disk solver.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prepare called")
	},
}

func init() {
	diskCmd.AddCommand(diskPrepareCmd)
}
