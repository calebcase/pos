package pos

import (
	"crypto/rand"
	"io"
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
	Claim        int64 // The amount of space in bytes for this puzzle.
	PRNG         PRNG  // The PRNG to use for generating data for the puzzle.
	IndexSize    int64 // The size in bytes of the solution indices.
	SolutionSize int64 // The size in bytes of the solution.
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

func (p *Puzzle) PreseedIndices(n int64, seed []byte) (indices []int64, err error) {
	indices, err = p.selectIndices(n-1, seed)
	if err != nil {
		return nil, err
	}

	indices = append(indices, p.Claim-1)

	return indices, nil
}

// Indices computes the offsets of the solution bytes. Read the byte at each
// offset to create a solution.
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

// A type implementing the Solver interface can be used to prepare and solve a
// given puzzle.
type Solver interface {
	Prepare(puzzle *Puzzle) error
	Solve(preseedIndices []int64, mask []byte) (solution []byte, err error)
}
