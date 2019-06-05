package pos

import "io"

type DiskSolver struct {
	out io.ReadWriteSeeker
}

var _ Solver = (*DiskSolver)(nil)

func NewDiskSolver(out io.ReadWriteSeeker) (*DiskSolver, error) {
	return &DiskSolver{
		out: out,
	}, nil
}

func (s *DiskSolver) Prepare(puzzle *Puzzle) (err error) {
	prng, err := puzzle.PRNG.Clone()
	if err != nil {
		return err
	}

	const lastSize = 1024
	last := make([]byte, lastSize, lastSize)

	for i := int64(0); i < puzzle.Claim; i += lastSize {
		_, err := io.ReadFull(prng, last)
		if err != nil {
			return err
		}

		_, err = s.out.Write(last)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DiskSolver) fromIndices(indices []int64) (value []byte, err error) {
	value = make([]byte, len(indices), len(indices))

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

		value[i] = tmp[0]
	}

	return value, nil
}

func (s *DiskSolver) Solve(puzzle *Puzzle, preseedIndices []int64, mask []byte) (solution []byte, err error) {
	var preseed []byte

	// First Pass: Read all preseed indices and construct the preseed.
	for i := int64(0); i < puzzle.PreseedRounds; i++ {
		preseed, err = s.fromIndices(preseedIndices)
		if err != nil {
			return nil, err
		}

		preseedIndices, err = puzzle.PreseedIndices(int64(len(preseedIndices)), preseed)
		if err != nil {
			return nil, err
		}

	}

	preseed, err = s.fromIndices(preseedIndices)
	if err != nil {
		return nil, err
	}

	// Second Pass: Read all the solution indices and construct the solution.
	solutionIndices, err := puzzle.SolutionIndices(preseed, mask)
	if err != nil {
		return nil, err
	}

	solution, err = s.fromIndices(solutionIndices)
	if err != nil {
		return nil, err
	}

	return solution, nil
}
