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

	claim := int64(1024 * 1024 * 1024 * 1)
	fmt.Printf("Claim: %d\n", claim)

	key, err := pos.NewRandomBytes(32)
	if err != nil {
		panic(err)
	}

	iv, err := pos.NewRandomBytes(16)
	if err != nil {
		panic(err)
	}

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

	prng, err := aesprng.New(key, iv)
	if err != nil {
		panic(err)
	}

	puzzle := &pos.Puzzle{
		Claim:        claim,
		PRNG:         prng,
		IndexSize:    64,
		SolutionSize: 16,
	}

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

	start = time.Now()

	streamSolution, err := streamSolver.Solve(preseedIndices, mask)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Stream Solution (%s):\n%x\n", time.Since(start), streamSolution)

	start = time.Now()

	diskSolution, err := diskSolver.Solve(preseedIndices, mask)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Disk Solution (%s):\n%x\n", time.Since(start), diskSolution)
}
