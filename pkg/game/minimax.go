package game

import (
	"math"
	"sort"
)

var MaxDepth = 13

type minimax struct {
	s side
}

func NewMinimax(s side) *minimax {
	return &minimax{
		s: s,
	}
}

func (m *minimax) UpdateSide(s side) {
	m.s = s
}

func (m *minimax) GetBestMove(b *board) (int, error) {
	max := int64(math.MinInt64)
	maxHole := 0
	heuristicValues, err := m.getPossibleMoves(b, m.s)
	if err != nil {
		return -1, err
	}
	sort.Sort(sort.Reverse(heuristicValues))
	alpha := int64(math.MinInt64)
	beta := int64(math.MaxInt64)
	//clone := b.clone()
	//nextSide, err := clone.makeMove(m.s, heuristicValues[0].hole)
	//if err != nil {
	//return -1, err
	//}
	//score := m.do(clone, nextSide, nextSide == m.s, MaxDepth, alpha, beta)
	//max = score
	//maxHole = heuristicValues[0].hole
	//alpha = max
	//heuristicValues = heuristicValues[1:]
	for _, node := range heuristicValues {
		clone := b.clone()
		nextSide, err := clone.makeMove(m.s, node.hole)
		if err != nil {
			return -1, err
		}
		score, err := m.do(clone, nextSide, nextSide == m.s, MaxDepth, alpha, beta)
		if err != nil {
			return -1, err
		}
		if score > max {
			max = score
			maxHole = node.hole
			if max > alpha {
				alpha = max
			}
		}
	}
	return maxHole, nil
}

func (m *minimax) do(b *board, s side, max bool, depth int, alpha int64, beta int64) (int64, error) {
	if depth == 0 || b.gameOver() {
		return m.boardHeuristic(b, s)
	}
	maxScore := int64(math.MinInt64)
	minScore := int64(math.MaxInt64)
	heuristicValues, err := m.getPossibleMoves(b, s)
	if err != nil {
		return 0, err
	}
	if max {
		sort.Sort(sort.Reverse(heuristicValues))
	} else {
		sort.Sort(heuristicValues)
	}
	for _, node := range heuristicValues {
		clone := b.clone()
		nextSide, err := clone.makeMove(s, node.hole)
		if err != nil {
			return 0, err
		}
		score, err := m.do(clone, nextSide, nextSide == m.s, depth-1, alpha, beta)
		if err != nil {
			return 0, err
		}
		if max && score > maxScore {
			maxScore = score
			if maxScore > alpha {
				alpha = maxScore
			}
		} else if !max && score < minScore {
			minScore = score
			if minScore < beta {
				beta = minScore
			}
		}
		if beta <= alpha {
			if max {
				return maxScore, nil
			}
			return minScore, nil
		}
	}
	if max {
		return maxScore, nil
	}
	return minScore, nil
}

func (m *minimax) getPossibleMoves(b *board, s side) (heuristicValues, error) {
	h := heuristicValues{}
	for i := 1; i <= b.getNoOfHoles(); i++ {
		if b.isLegal(s, i) {
			heuristic, err := m.moveHeuristic(b, s, i)
			if err != nil {
				return heuristicValues{}, err
			}
			h = append(h, heuristicValue{hole: i, heuristic: heuristic})
		}
	}
	return h, nil
}

func (m *minimax) moveHeuristic(b *board, s side, hole int) (int64, error) {
	score := int64(0)
	clone := b.clone()
	nextSide, err := clone.makeMove(s, hole)
	if err != nil {
		return 0, err
	}
	if nextSide == m.s {
		score += 10
	} else {
		score -= 10
	}
	h, err := m.boardHeuristic(clone, nextSide)
	if err != nil {
		return 0, err
	}
	score += h
	return score, nil
}

func (m *minimax) boardHeuristic(b *board, s side) (int64, error) {
	score := int64(b.state[index(m.s)][0] - b.state[index(m.s.Opposite())][0])
	seedsOnSide := int64(0)
	for i := 5; i <= b.getNoOfHoles(); i++ {
		seeds, err := b.getSeeds(m.s, i)
		if err != nil {
			return 0, err
		}
		seedsOnSide += int64(seeds)
	}
	score += seedsOnSide / 8
	if b.state[index(m.s)][0] > 49 {
		score += 100
	} else if b.state[index(m.s.Opposite())][0] > 49 {
		score -= 100
	}
	maxSteal := int64(0)
	for i := 1; i <= b.getNoOfHoles(); i++ {
		seeds, err := b.getSeeds(s, i)
		if err != nil {
			return 0, err
		}
		land := (i + seeds) % 15
		if land > 0 && land <= 7 {
			landSeeds, err := b.getSeeds(s, land)
			if err != nil {
				return 0, err
			}
			if seeds > 0 && seeds <= 15 && (landSeeds == 0 || (seeds == 15 && land == i)) {
				seedsOpp, err := b.getSeedsOp(s, land)
				if err != nil {
					return 0, err
				}
				if seedsOpp > 0 {
					if maxSteal < int64(seedsOpp+1) {
						maxSteal = int64(seedsOpp + 1)
					}
				}
			}
		}
	}
	if s == m.s {
		score += maxSteal
	} else {
		score -= maxSteal
	}
	return score, nil
}

type heuristicValues []heuristicValue

func (h heuristicValues) Len() int {
	return len(h)
}

func (h heuristicValues) Less(i, j int) bool {
	return h[i].heuristic < h[j].heuristic
}

func (h heuristicValues) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

type heuristicValue struct {
	hole      int
	heuristic int64
}
