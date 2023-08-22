package usecase

import (
	"reversi/src/domain"
	"reversi/src/utility/time"
)

type SoloReversi struct {
	*BaseReversi
	ra ReversiAlgorithm
}

func (s *SoloReversi) PostPut() {
	r := s.BaseReversi

	cpStone := r.currentPlayerStone
	cells := r.PutableFieldCells(cpStone)

	time.Wait(1, 2)
	s.ra.put(s, cpStone, cells)
}

func (s *SoloReversi) PostPass() {
	s.PostPut()
}

func (s *SoloReversi) IsSoloPlay() bool {
	return true
}

type ReversiAlgorithm interface {
	put(s *SoloReversi, cpStone domain.Stone, cells []domain.PutableFieldCell)
}

type SimpleReversiAlgolithm struct{}

func (ra *SimpleReversiAlgolithm) put(s *SoloReversi, cpStone domain.Stone, cells []domain.PutableFieldCell) {
	if len(cells) == 0 {
		s.Pass()
		return
	}

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

	s.Put(best)
}
