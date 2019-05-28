package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "stream commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stream called")
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)
}
