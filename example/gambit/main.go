// This example runs the gambit of features provided by the pos library. It
// will create a stream and disk solver and go through all phases for a proof
// of space.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/calebcase/pos"
	"github.com/calebcase/pos/lib/aesprng"
)

func main() {
	var start time.Time

	// Prepare a claim of 1 GiB.
	claim := int64(1024 * 1024 * 1024 * 1)
	fmt.Printf("Claim: %d\n", claim)

	// Generate a random key and iv for the AES PRNG.
	key, err := pos.NewRandomBytes(32)
	if err != nil {
		panic(err)
	}

	iv, err := pos.NewRandomBytes(16)
	if err != nil {
		panic(err)
	}

	prng, err := aesprng.New(key, iv)
	if err != nil {
		panic(err)
	}

	// Initialize the puzzle with some reasonable defaults.
	puzzle := &pos.Puzzle{
		Claim:         claim,
		PRNG:          prng,
		PreseedRounds: 10,
		IndexSize:     64,
		SolutionSize:  16,
	}

	// Create the stream and disk solvers. We will compare their results
	// and runtime.
	streamSolver, err := pos.NewStreamSolver()
	if err != nil {
		panic(err)
	}

	file, err := os.Create("output")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	diskSolver, err := pos.NewDiskSolver(file)
	if err != nil {
		panic(err)
	}

	// Time the stream and disk solver's prepare phase.
	start = time.Now()

	err = streamSolver.Prepare(puzzle)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Stream Solver Prepared (%s)\n", time.Since(start))

	start = time.Now()

	err = diskSolver.Prepare(puzzle)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Disk Solver Prepared (%s)\n", time.Since(start))

	file.Sync()

	// Select preseed indices and a mask for the solution phase.
	preseedIdxSeed, err := pos.NewRandomBytes(len(key) + len(iv))
	if err != nil {
		panic(err)
	}

	preseedIndices, err := puzzle.PreseedIndices(int64(len(preseedIdxSeed)), preseedIdxSeed)
	if err != nil {
		panic(err)
	}

	mask, err := pos.NewRandomBytes(len(key) + len(iv))
	if err != nil {
		panic(err)
	}

	// Time the stream and disk solver's solve phase.
	start = time.Now()

	streamSolution, err := streamSolver.Solve(puzzle, preseedIndices, mask)
	if err != nil {
		panic(err)
	}

	streamSolutionTime := time.Since(start)

	fmt.Printf("Stream Solution (%s):\n%x\n", streamSolutionTime, streamSolution)

	start = time.Now()

	diskSolution, err := diskSolver.Solve(puzzle, preseedIndices, mask)
	if err != nil {
		panic(err)
	}

	diskSolutionTime := time.Since(start)

	fmt.Printf("Disk Solution (%s):\n%x\n", diskSolutionTime, diskSolution)

	fmt.Printf("Solution Ratio: %f\n", float64(streamSolutionTime)/float64(diskSolutionTime))
}
