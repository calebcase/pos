package cmd

import (
	"encoding/json"
	"os"

	"github.com/calebcase/pos"
	"github.com/spf13/cobra"
)

var diskPrepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Prepare a disk solver",
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

		path, err := cmd.Flags().GetString("image")
		if err != nil {
			panic(err)
		}

		image, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer image.Close()

		diskSolver, err := pos.NewDiskSolver(image)
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

		err = diskSolver.Prepare(puz)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	diskCmd.AddCommand(diskPrepareCmd)
}
