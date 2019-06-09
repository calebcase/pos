package cmd

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/calebcase/pos"
	"github.com/spf13/cobra"
)

var puzzleMaskCmd = &cobra.Command{
	Use:   "mask",
	Short: "Create a mask for a puzzle",
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

		mask, err := pos.NewRandomBytes(len(p.PRNG.GetSeed()))
		if err != nil {
			panic(err)
		}

		_, err = os.Stdout.Write([]byte(base64.StdEncoding.EncodeToString(mask) + "\n"))
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	puzzleCmd.AddCommand(puzzleMaskCmd)

	puzzleMaskCmd.PersistentFlags().StringP("puzzle", "p", "", "Path to a puzzle config")
}
