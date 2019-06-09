package cmd

import (
	"github.com/calebcase/pos"
	"github.com/calebcase/pos/lib/aesprng"
	"github.com/spf13/cobra"
)

var puzzleCmd = &cobra.Command{
	Use:   "puzzle",
	Short: "Puzzle commands",
}

type puzzle struct {
	pos.Puzzle

	PRNG *aesprng.State `json:"prng"`
}

func init() {
	rootCmd.AddCommand(puzzleCmd)
}
