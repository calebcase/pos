package cmd

import (
	"encoding/json"
	"os"

	"github.com/calebcase/pos"
	"github.com/spf13/cobra"
)

var puzzlePreseedIndicesCmd = &cobra.Command{
	Use:   "preseed-indices",
	Short: "Compute preseed indices for a puzzle.",
	Run: func(cmd *cobra.Command, args []string) {
		input := os.Stdin
		if cmd.Flags().Changed("puzzle") {
			path, err := cmd.Flags().GetString("puzzle")
			if err != nil {
				panic(err)
			}

			if path != "-" {
				input, err = os.Open(path)
				if err != nil {
					panic(err)
				}
				defer input.Close()
			}
		}

		var p puzzle

		err := json.NewDecoder(input).Decode(&p)
		if err != nil {
			panic(err)
		}

		p.Puzzle.PRNG = p.PRNG

		seed, err := pos.NewRandomBytes(len(p.PRNG.GetSeed()))
		if err != nil {
			panic(err)
		}

		indices, err := p.PreseedIndices(int64(len(seed)), seed)
		if err != nil {
			panic(err)
		}

		err = json.NewEncoder(os.Stdout).Encode(indices)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	puzzleCmd.AddCommand(puzzlePreseedIndicesCmd)

	puzzlePreseedIndicesCmd.PersistentFlags().StringP("puzzle", "p", "", "Path to a puzzle config")
}
