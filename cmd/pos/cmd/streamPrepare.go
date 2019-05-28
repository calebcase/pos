package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var streamPrepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Prepare a stream solver.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prepare called")
	},
}

func init() {
	streamCmd.AddCommand(streamPrepareCmd)
}
