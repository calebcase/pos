package pos

import "io"

type StreamSolver struct {
	puzzle *Puzzle
}

var _ Solver = (*StreamSolver)(nil)

func NewStreamSolver() (*StreamSolver, error) {
	return &StreamSolver{}, nil
}

func (s *StreamSolver) Prepare(puzzle *Puzzle) (err error) {
	s.puzzle = puzzle

	return nil
}

func (s *StreamSolver) fromIndices(indices []int64) (value []byte, err error) {
	prng, err := s.puzzle.PRNG.Clone()
	if err != nil {
		return nil, err
	}

	mapper := make(map[int64]byte)

	const lastSize = 1024
	last := make([]byte, lastSize, lastSize)

	for i := int64(0); i < s.puzzle.Claim; i += lastSize {
		_, err := io.ReadFull(prng, last)
		if err != nil {
			return nil, err
		}

		for _, idx := range indices {
			if idx >= i && idx < i+lastSize {
				mapper[idx] = last[idx-i]
			}
		}
	}

	value = make([]byte, len(indices), len(indices))

	for i, index := range indices {
		value[i] = mapper[index]
	}

	return value, nil
}

func (s *StreamSolver) Solve(preseedIndices []int64, mask []byte) (solution []byte, err error) {
	// First Pass: Read all preseed indices and construct the preseed.
	preseed, err := s.fromIndices(preseedIndices)
	if err != nil {
		return nil, err
	}

	// Second Pass: Read all the solution indices and construct the solution.
	solutionIndices, err := s.puzzle.SolutionIndices(preseed, mask)
	if err != nil {
		return nil, err
	}

	solution, err = s.fromIndices(solutionIndices)
	if err != nil {
		return nil, err
	}

	return solution, nil
}
