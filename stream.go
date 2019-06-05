package pos

import "io"

type StreamSolver struct{}

var _ Solver = (*StreamSolver)(nil)

func NewStreamSolver() (*StreamSolver, error) {
	return &StreamSolver{}, nil
}

func (s *StreamSolver) Prepare(puzzle *Puzzle) (err error) {
	return nil
}

func (s *StreamSolver) fromIndices(puzzle *Puzzle, indices []int64) (value []byte, err error) {
	prng, err := puzzle.PRNG.Clone()
	if err != nil {
		return nil, err
	}

	mapper := make(map[int64]byte)

	const lastSize = 1024
	last := make([]byte, lastSize, lastSize)

	for i := int64(0); i < puzzle.Claim; i += lastSize {
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

func (s *StreamSolver) Solve(puzzle *Puzzle, preseedIndices []int64, mask []byte) (solution []byte, err error) {
	var preseed []byte

	// First Pass: Read all preseed indices and construct the preseed.
	for i := int64(0); i < puzzle.PreseedRounds; i++ {
		preseed, err = s.fromIndices(puzzle, preseedIndices)
		if err != nil {
			return nil, err
		}

		preseedIndices, err = puzzle.PreseedIndices(int64(len(preseedIndices)), preseed)
		if err != nil {
			return nil, err
		}
	}

	preseed, err = s.fromIndices(puzzle, preseedIndices)
	if err != nil {
		return nil, err
	}

	// Second Pass: Read all the solution indices and construct the solution.
	solutionIndices, err := puzzle.SolutionIndices(preseed, mask)
	if err != nil {
		return nil, err
	}

	solution, err = s.fromIndices(puzzle, solutionIndices)
	if err != nil {
		return nil, err
	}

	return solution, nil
}
