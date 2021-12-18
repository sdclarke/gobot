package game

import "errors"

type board struct {
	holes int
	state [][]int
}

func index(s side) int {
	if s == North {
		return 0
	}
	return 1
}

func NewBoard(holes, seeds int) *board {
	b := &board{
		holes: holes,
		state: make([][]int, 2),
	}
	for i := range b.state {
		b.state[i] = make([]int, holes+1)
		for j := 1; j <= holes; j++ {
			b.state[i][j] = seeds
		}
	}
	return b
}

func (b *board) getNoOfHoles() int {
	return b.holes
}

func (b *board) setSeeds(s side, hole, seeds int) error {
	if hole < 1 || hole > b.holes {
		return errors.New("Invalid hole number")
	}
	if seeds < 0 {
		return errors.New("Expected non-negative number of seeds")
	}
	b.state[index(s)][hole] = seeds
	return nil
}

func (b *board) setSeedsOp(s side, hole, seeds int) error {
	if hole < 1 || hole > b.holes {
		return errors.New("Invalid hole number")
	}
	if seeds < 0 {
		return errors.New("Expected non-negative number of seeds")
	}
	b.state[index(s.Opposite())][b.holes-hole+1] = seeds
	return nil
}

func (b *board) getSeeds(s side, hole int) (int, error) {
	if hole < 1 || hole > b.holes {
		return -1, errors.New("Invalid hole number")
	}
	return b.state[index(s)][hole], nil
}

func (b *board) getSeedsOp(s side, hole int) (int, error) {
	if hole < 1 || hole > b.holes {
		return -1, errors.New("Invalid hole number")
	}
	return b.state[index(s.Opposite())][b.holes-hole+1], nil
}

func (b *board) setSeedsInStore(s side, seeds int) error {
	if seeds < 0 {
		return errors.New("Expected non-negative number of seeds")
	}
	b.state[index(s)][0] = seeds
	return nil
}

func (b *board) makeMove(s side, hole int) (side, error) {
	seedsToSow, err := b.getSeeds(s, hole)
	if err != nil {
		return None, err
	}
	err = b.setSeeds(s, hole, 0)
	if err != nil {
		return None, err
	}
	receivingPits := 2*b.holes + 1
	rounds := seedsToSow / receivingPits
	extra := seedsToSow % receivingPits
	if rounds > 0 {
		for i := 1; i <= b.holes; i++ {
			b.state[index(South)][i] += rounds
			b.state[index(North)][i] += rounds
		}
		b.state[index(s)][0] += rounds
	}

	sowSide := s
	sowHole := hole
	for ; extra > 0; extra-- {
		sowHole++
		if sowHole == 1 {
			sowSide = sowSide.Opposite()
		}
		if sowHole > b.holes {
			if sowSide == s {
				sowHole = 0
				b.state[index(s)][0] += 1
				continue
			} else {
				sowSide = sowSide.Opposite()
				sowHole = 1
			}
		}
		b.state[index(sowSide)][sowHole] += 1
	}

	if sowSide == s && sowHole > 0 {
		sowHoleSeeds, err := b.getSeeds(sowSide, sowHole)
		if err != nil {
			return None, err
		}
		sowHoleSeedsOp, err := b.getSeedsOp(sowSide, sowHole)
		if err != nil {
			return None, err
		}
		if sowHoleSeeds == 1 && sowHoleSeedsOp > 0 {
			b.state[index(s)][0] += 1 + sowHoleSeedsOp
			err = b.setSeeds(s, sowHole, 0)
			if err != nil {
				return None, err
			}
			err = b.setSeedsOp(s, sowHole, 0)
			if err != nil {
				return None, err
			}
		}
	}

	finishedSide := None
	if b.holesEmpty(s) {
		finishedSide = s
	} else if b.holesEmpty(s.Opposite()) {
		finishedSide = s.Opposite()
	}
	if finishedSide != None {
		seeds := 0
		collectingSide := finishedSide.Opposite()
		for i := 1; i <= b.holes; i++ {
			seeds += b.state[index(collectingSide)][i]
			err = b.setSeeds(collectingSide, i, 0)
			if err != nil {
				return None, err
			}
		}
		b.state[index(collectingSide)][0] += seeds
	}
	if sowHole == 0 {
		return s, nil
	}
	return s.Opposite(), nil
}

func (b *board) gameOver() bool {
	return b.holesEmpty(South) || b.holesEmpty(North)
}

func (b *board) holesEmpty(s side) bool {
	for i := 1; i <= b.holes; i++ {
		if b.state[index(s)][i] != 0 {
			return false
		}
	}
	return true
}

func (b *board) isLegal(s side, hole int) bool {
	if hole < 1 || hole > b.holes {
		return false
	}
	if b.state[index(s)][hole] < 1 {
		return false
	}
	return true
}

func (b *board) clone() *board {
	state := make([][]int, len(b.state))
	for i := range b.state {
		state[i] = make([]int, len(b.state[i]))
		copy(state[i], b.state[i])
	}
	return &board{
		holes: b.holes,
		state: state,
	}
}
