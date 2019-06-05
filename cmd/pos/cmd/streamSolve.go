package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/calebcase/pos"
	"github.com/spf13/cobra"
)

var streamSolveCmd = &cobra.Command{
	Use:   "solve",
	Short: "Solve a puzzle with a stream solver.",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

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

		dec := json.NewDecoder(input)

		var p puzzle

		err = dec.Decode(&p)
		if err != nil {
			panic(err)
		}

		streamSolver, err := pos.NewStreamSolver()
		if err != nil {
			panic(err)
		}

		puz := &pos.Puzzle{
			Claim:         p.Claim,
			PRNG:          p.PRNG,
			PreseedRounds: p.PreseedRounds,
			IndexSize:     p.IndexSize,
			SolutionSize:  p.SolutionSize,
		}

		preseedIndices, err := cmd.Flags().GetInt64Slice("preseed-indices")
		if err != nil {
			panic(err)
		}

		mask, err := cmd.Flags().GetBytesBase64("mask")
		if err != nil {
			panic(err)
		}

		solution, err := streamSolver.Solve(puz, preseedIndices, mask)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Solution: %x\n", solution)
	},
}

func init() {
	streamCmd.AddCommand(streamSolveCmd)

	streamSolveCmd.PersistentFlags().Int64Slice("preseed-indices", []int64{}, "A list of preseed indices")
	cobra.MarkFlagRequired(streamSolveCmd.PersistentFlags(), "preseed-indices")

	streamSolveCmd.PersistentFlags().BytesBase64("mask", []byte{}, "A base64 encoded mask")
	cobra.MarkFlagRequired(streamSolveCmd.PersistentFlags(), "mask")
}
