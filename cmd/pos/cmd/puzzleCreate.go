package cmd

import (
	"encoding/json"
	"os"

	"github.com/calebcase/pos"
	"github.com/calebcase/pos/lib/aesprng"
	"github.com/spf13/cobra"
)

var puzzleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new puzzle",
	Run: func(cmd *cobra.Command, args []string) {
		claim, err := cmd.Flags().GetInt64("claim")
		if err != nil {
			panic(err)
		}

		seed, err := cmd.Flags().GetBytesBase64("seed")
		if err != nil {
			panic(err)
		}

		key, iv, err := aesprng.SplitSeed(seed)
		if err != nil {
			panic(err)
		}

		prng, err := aesprng.New(key, iv)
		if err != nil {
			panic(err)
		}

		indexSize, err := cmd.Flags().GetInt64("index-size")
		if err != nil {
			panic(err)
		}

		solutionSize, err := cmd.Flags().GetInt64("solution-size")
		if err != nil {
			panic(err)
		}

		var preseedRounds int64

		if cmd.Flags().Changed("preseed-rounds") {
			preseedRounds, err = cmd.Flags().GetInt64("preseed-rounds")
			if err != nil {
				panic(err)
			}
		} else {
			rate, err := cmd.Flags().GetFloat64("pr-est-rate")
			if err != nil {
				panic(err)
			}

			scale, err := cmd.Flags().GetFloat64("pr-est-scale")
			if err != nil {
				panic(err)
			}

			preseedRounds = pos.EstimatePreseedRounds(claim, rate, scale)
		}

		puzzle := &pos.Puzzle{
			Claim:         claim,
			PRNG:          prng,
			PreseedRounds: preseedRounds,
			IndexSize:     indexSize,
			SolutionSize:  solutionSize,
		}

		err = json.NewEncoder(os.Stdout).Encode(puzzle)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	puzzleCmd.AddCommand(puzzleCreateCmd)

	puzzleCreateCmd.PersistentFlags().Int64P("claim", "c", 0, "Size of the claimed storage (bytes)")
	cobra.MarkFlagRequired(puzzleCreateCmd.PersistentFlags(), "claim")

	defaultSeed, _ := pos.NewRandomBytes(32 + 16)
	puzzleCreateCmd.PersistentFlags().BytesBase64("seed", defaultSeed, "A base64 encoded seed")

	puzzleCreateCmd.PersistentFlags().Int64("index-size", 64, "Size of the index (bytes)")
	puzzleCreateCmd.PersistentFlags().Int64("solution-size", 10, "Size of the solution (bytes)")

	puzzleCreateCmd.PersistentFlags().Int64("preseed-rounds", 0, "Number of preseed rounds")
	puzzleCreateCmd.PersistentFlags().Float64("pr-est-rate", 1024*1024*1024*10, "Rate of PRNG generation (bytes per second)")
	puzzleCreateCmd.PersistentFlags().Float64("pr-est-scale", 2, "Desired time scale (seconds)")
}
