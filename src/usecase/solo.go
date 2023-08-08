package usecase

import (
	"reversi/src/domain"
)

type SoloReversiHandler struct {
	ra ReversiAlgorithm
}

func (s *SoloReversiHandler) postPut(r *Reversi) domain.Stone {
	cpStone := domain.SwitchStone(r.currentPlayerStone)
	cells := r.PutableFieldCells(cpStone)
	s.ra.put(r, cpStone, cells)
	return r.currentPlayerStone
}

type ReversiAlgorithm interface {
	put(r *Reversi, cpStone domain.Stone, cells []domain.PutableFieldCell)
}

type SimpleReversiAlgolithm struct{}

func (s *SimpleReversiAlgolithm) put(r *Reversi, cpStone domain.Stone, cells []domain.PutableFieldCell) {
	if len(cells) == 0 {
		return
	}

	// fmt.Printf("%#+v", cells)

	weight := [][]int{
		{30, -10, 2, 1, 1, 2, -10, 30},
		{-10, -20, -3, -3, -3, -3, -20, -10},
		{2, -3, 2, 0, 0, 2, -3, 2},
		{1, -3, 0, 0, 0, 0, -3, 1},
		{1, -3, 0, 0, 0, 0, -3, 1},
		{2, -3, 2, 0, 0, 2, -3, 2},
		{-10, -20, -3, -3, -3, -3, -20, -10},
		{30, -10, 2, 1, 1, 2, -10, 30},
	}
	best := cells[0]
	for _, cell := range cells[1:] {
		if score := weight[cell.Pos.X][cell.Pos.Y]; score > weight[best.Pos.X][best.Pos.Y] {
			best = cell
		}
	}

	r.put(best)
}
