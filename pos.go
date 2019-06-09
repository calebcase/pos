package pos

import (
	"crypto/rand"
	"io"
	"math"
	"math/big"
)

// A type implementing the PRNG interface can be used to generate pseudo random
// numbers.
type PRNG interface {
	io.Reader

	// Create a new PRNG of the same type initialized with the given seed.
	New(seed []byte) (prng PRNG, err error)

	// Create a new PRNG of the same type initialized with the original seed.
	Clone() (prng PRNG, err error)

	// Get the seed used to initialize this PRNG.
	GetSeed() []byte
}

// NewRandomBytes returns random bytes of the given size.
func NewRandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Puzzle contains the parameters for preparing and solving a proof of space
// puzzle.
type Puzzle struct {
	Claim         int64 `json:"claim"`          // The amount of space in bytes for this puzzle.
	PRNG          PRNG  `json:"prng"`           // The PRNG to use for generating data for the puzzle.
	PreseedRounds int64 `json:"preseed_rounds"` // The number of rounds to compute during the preseed phase.
	IndexSize     int64 `json:"index_size"`     // The size in bytes of the solution indices.
	SolutionSize  int64 `json:"solution_size"`  // The size in bytes of the solution.
}

func (p *Puzzle) selectIndices(n int64, seed []byte) (indices []int64, err error) {
	prng, err := p.PRNG.New(seed)
	if err != nil {
		return nil, err
	}

	index := make([]byte, p.IndexSize, p.IndexSize)

	base := big.NewInt(p.Claim)
	ith := &big.Int{}

	indices = make([]int64, 0, n)

	for j := int64(0); j < p.Claim && int64(len(indices)) < n; j += p.IndexSize {
		_, err := io.ReadFull(prng, index)
		if err != nil {
			return nil, err
		}

		ith.SetBytes(index)
		ith.Mod(ith, base)

		indices = append(indices, ith.Int64())
	}

	return indices, nil
}

// PreseedIndices computes the offsets of the preseed bytes. Read the byte at
// each offset to create a preseed.
func (p *Puzzle) PreseedIndices(n int64, seed []byte) (indices []int64, err error) {
	indices, err = p.selectIndices(n-1, seed)
	if err != nil {
		return nil, err
	}

	indices = append(indices, p.Claim-1)

	return indices, nil
}

// SolutionIndices computes the offsets of the solution bytes. Read the byte at
// each offset to create a solution.
func (p *Puzzle) SolutionIndices(preseed, mask []byte) (indices []int64, err error) {
	seed := make([]byte, len(mask), len(mask))
	for i, _ := range preseed {
		seed[i] = preseed[i] ^ mask[i]
	}

	indices, err = p.selectIndices(p.SolutionSize, seed)
	if err != nil {
		return nil, err
	}

	return indices, nil
}

// EstimatePreseedRounds estimates the number of preseed rounds needed for the
// given claim assuming a given hashing rate in bytes per second and the
// desired minimum timescale in seconds.
//
// PRNG rate should be set to the fastest rate you believe can be acheived for
// your threat profile. For example, if your threat profile includes hardware
// accelerated PRNGs, and you have reviewed the hardware currently available to
// find the current state of the art is approximately 1 GiB / second, then a
// rate of 10x that may give you sufficient head room.
//
// The minimum timescale is meant to tune the preseed rounds such that it takes
// at least scale seconds long to generate one claim's worth of PRNG bytes.
//
// NOTE: The preseed rounds do not increase the differential between stream and
// disk solving. The preseed rounds are only meant to provide a scaling factor
// to the overall process to ensure that the differential is detectable. For
// example, if you are operating this over a network with average latency of 1
// second with a standard deviation of 0.5 seconds you may wish to set the
// scale to 2 * 1.5 = 3 seconds in order to ensure it is possible to
// differentiate between a stream and disk solution (without being washed out
// by the variance in the network itself).
func EstimatePreseedRounds(claim int64, rate, scale float64) int64 {
	unscaled := float64(claim) / rate

	if unscaled > scale {
		return 0
	}

	return int64(math.Ceil(scale / unscaled))
}

// A type implementing the Solver interface can be used to prepare and solve a
// given puzzle.
type Solver interface {
	Prepare(puzzle *Puzzle) error
	Solve(puzzle *Puzzle, preseedIndices []int64, mask []byte) (solution []byte, err error)
}
