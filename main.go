package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"
)

// NewRandomBytes returns random bytes of length size.
func NewRandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// A type implementing the PRNG interface can be used to generate pseudo random
// numbers.
type PRNG interface {
	io.Reader

	// Create a new PRNG of the same type. If seed is not nil, then the new
	// PRNG will use it as the new seed. Otherwise the seed will be
	// initialized to the current PRNGs value.
	New(seed []byte) (prng PRNG, err error)
}

type AESPRNG struct {
	key []byte
	iv  []byte

	mode cipher.BlockMode
	zero []byte
}

var _ PRNG = (*AESPRNG)(nil)

func NewAESPRNG(key, iv []byte) (prng *AESPRNG, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &AESPRNG{
		key: append([]byte(nil), key...),
		iv:  append([]byte(nil), iv...),

		mode: cipher.NewCBCEncrypter(block, iv),
		zero: make([]byte, 1024, 1024),
	}, nil
}

func (prng *AESPRNG) Read(b []byte) (n int, err error) {
	if len(prng.zero) < len(b) {
		prng.zero = make([]byte, len(b), len(b))
	}

	prng.mode.CryptBlocks(b, prng.zero[:len(b)])

	return len(b), nil
}

func (prng *AESPRNG) New(seed []byte) (clone PRNG, err error) {
	key := prng.key
	iv := prng.iv

	if seed != nil {
		key = seed[:32]
		iv = seed[32:]
	}

	return NewAESPRNG(key, iv)
}

type Puzzle struct {
	Type string

	Claim int64

	Seed []byte

	IndexSize int64

	SolutionSize int64
}

func NewPuzzle(claim int64) (*Puzzle, error) {
	seed, err := NewRandomBytes(32 + 16)
	if err != nil {
		return nil, err
	}

	return &Puzzle{
		Type: "aes-cbc-256",

		Claim: claim,

		Seed: seed,

		IndexSize: 1024,

		SolutionSize: 16,
	}, nil
}

func (p *Puzzle) Indices(prng PRNG, last, mask []byte) (indices []int64, err error) {
	index := make([]byte, p.IndexSize, p.IndexSize)

	base := big.NewInt(p.Claim)

	seed := make([]byte, len(mask), len(mask))
	for i, _ := range last {
		seed[i] = last[i] ^ mask[i]
	}

	iprng, err := prng.New(seed)
	if err != nil {
		return nil, err
	}

	indices = []int64{}
	for j := int64(0); j < p.Claim && int64(len(indices)) < p.SolutionSize; j += p.IndexSize {
		_, err := io.ReadFull(iprng, index)
		if err != nil {
			return nil, err
		}

		ith := &big.Int{}
		ith.SetBytes(index)
		ith.Mod(ith, base)

		indices = append(indices, ith.Int64())
	}

	return indices, nil
}

type Solver interface {
	Prepare(puzzle *Puzzle) error
	Solve(mask []byte) (solution []byte, err error)
}

type StreamSolver struct {
	puzzle *Puzzle
	prng   PRNG

	last []byte
}

var _ Solver = (*StreamSolver)(nil)

func NewStreamSolver() (*StreamSolver, error) {
	return &StreamSolver{}, nil
}

func (s *StreamSolver) Prepare(puzzle *Puzzle) (err error) {
	s.puzzle = puzzle

	// Initialize a PRNG for the given puzzle type.
	switch s.puzzle.Type {
	case "aes-cbc-256":
		// AES CBC 256 key size is 32 and iv is 16.
		keySize := 32

		// Seed is packed as key + iv.
		key := s.puzzle.Seed[:keySize]
		iv := s.puzzle.Seed[keySize:]

		s.prng, err = NewAESPRNG(key, iv)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported puzzle type: %s", s.puzzle.Type)
	}

	// Our first pass is just to get the last block so that we can apply
	// the mask to it.

	// Create a buffer large enough to contain at least len(Seed) worth of
	// bytes.
	lastSize := 1024
	for lastSize < len(s.puzzle.Seed) {
		lastSize = lastSize * 2
	}

	last := make([]byte, lastSize)

	for i := int64(0); i < s.puzzle.Claim; i += int64(lastSize) {
		_, err := io.ReadFull(s.prng, last)
		if err != nil {
			return err
		}
	}

	s.last = last[len(last)-len(s.puzzle.Seed):]

	return nil
}

func (s *StreamSolver) Solve(mask []byte) (solution []byte, err error) {
	if mask == nil || len(mask) != len(s.puzzle.Seed) {
		return nil, fmt.Errorf("Invalid mask size %d. Expected mask size to match seed size of %d.", len(mask), len(s.puzzle.Seed))
	}

	if s.last == nil || len(s.last) != len(s.puzzle.Seed) {
		return nil, fmt.Errorf("Internal state is inconsistent. Prepare(puzzle) may not have been run.")
	}

	// Retrieve the indices we need to look for from the indices PRNG.
	iprng, err := s.prng.New(nil)
	if err != nil {
		return nil, err
	}

	indices, err := s.puzzle.Indices(iprng, s.last, mask)
	if err != nil {
		return nil, err
	}

	// Reset the PRNG for the 2nd pass. This time we are looking for the
	// indices we have identified.
	s.prng, err = s.prng.New(nil)
	if err != nil {
		return nil, err
	}

	// Create a buffer large enough to contain at least len(Seed) worth of
	// bytes.
	lastSize := 1024
	for lastSize < len(s.puzzle.Seed) {
		lastSize = lastSize * 2
	}

	last := make([]byte, lastSize)

	solution = make([]byte, s.puzzle.SolutionSize, s.puzzle.SolutionSize)

	mapper := make(map[int64]byte)

	for i := int64(0); i < s.puzzle.Claim; i += int64(lastSize) {
		_, err := io.ReadFull(s.prng, last)
		if err != nil {
			return nil, err
		}

		for _, idx := range indices {
			if idx >= i && idx < i+int64(lastSize) {
				mapper[idx] = last[idx-i]
			}
		}
	}

	for i, index := range indices {
		solution[i] = mapper[index]
	}

	return solution, nil
}

type DiskSolver struct {
	out io.ReadWriteSeeker

	puzzle *Puzzle
	prng   PRNG

	last []byte
}

var _ Solver = (*DiskSolver)(nil)

func NewDiskSolver(out io.ReadWriteSeeker) (*DiskSolver, error) {
	return &DiskSolver{
		out: out,
	}, nil
}

func (s *DiskSolver) Prepare(puzzle *Puzzle) (err error) {
	s.puzzle = puzzle

	// Initialize a PRNG for the given puzzle type.
	switch s.puzzle.Type {
	case "aes-cbc-256":
		// AES CBC 256 key size is 32 and iv is 16.
		keySize := 32

		// Seed is packed as key + iv.
		key := s.puzzle.Seed[:keySize]
		iv := s.puzzle.Seed[keySize:]

		s.prng, err = NewAESPRNG(key, iv)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported puzzle type: %s", s.puzzle.Type)
	}

	// Create a buffer large enough to contain at least len(Seed) worth of
	// bytes.
	lastSize := 1024
	for lastSize < len(s.puzzle.Seed) {
		lastSize = lastSize * 2
	}

	last := make([]byte, lastSize)

	for i := int64(0); i < s.puzzle.Claim; i += int64(lastSize) {
		_, err := io.ReadFull(s.prng, last)
		if err != nil {
			return err
		}

		_, err = s.out.Write(last)
		if err != nil {
			panic(err)
		}
	}

	s.last = last[len(last)-len(s.puzzle.Seed):]

	return nil
}

func (s *DiskSolver) Solve(mask []byte) (solution []byte, err error) {
	solution = make([]byte, s.puzzle.SolutionSize, s.puzzle.SolutionSize)

	iprng, err := s.prng.New(nil)
	if err != nil {
		return nil, err
	}

	indices, err := s.puzzle.Indices(iprng, s.last, mask)
	if err != nil {
		return nil, err
	}

	tmp := make([]byte, 1, 1)

	for i, index := range indices {
		_, err := s.out.Seek(index, io.SeekStart)
		if err != nil {
			return nil, err
		}

		_, err = io.ReadFull(s.out, tmp)
		if err != nil {
			return nil, err
		}

		solution[i] = tmp[0]
	}

	return solution, nil
}

func main() {
	var start time.Time

	claim := int64(1024 * 1024 * 1024 * 1)
	fmt.Printf("Claim: %d\n", claim)

	key, err := NewRandomBytes(32)
	if err != nil {
		panic(err)
	}

	iv, err := NewRandomBytes(16)
	if err != nil {
		panic(err)
	}

	streamSolver, err := NewStreamSolver()
	if err != nil {
		panic(err)
	}

	of, err := os.Create("output")
	if err != nil {
		panic(err)
	}
	defer of.Close()

	diskSolver, err := NewDiskSolver(of)
	if err != nil {
		panic(err)
	}

	puzzle, err := NewPuzzle(claim)
	if err != nil {
		panic(err)
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

	of.Sync()

	mask, err := NewRandomBytes(len(key) + len(iv))
	if err != nil {
		panic(err)
	}

	start = time.Now()

	streamSolution, err := streamSolver.Solve(mask)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Stream Solution (%s):\n%x\n", time.Since(start), streamSolution)

	start = time.Now()

	diskSolution, err := diskSolver.Solve(mask)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Disk Solution (%s):\n%x\n", time.Since(start), diskSolution)
}
