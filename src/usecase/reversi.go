package usecase

import (
	"fmt"
	"reversi/src/domain"
	"reversi/src/utility/slices"
)

type Reversi interface {
	PutableFieldCells(field domain.Field, playerStone domain.Stone) []domain.PutableFieldCell
	Put(field domain.Field, playerStone domain.Stone) domain.PutableFieldCell
	GetScore(field domain.Field) domain.Score
	IsFinished(field domain.Field) bool
	IsSoloPlay() bool
}

func NewReversi(solo, duel bool) Reversi {
	base := &BaseReversi{}
	if solo {
		return &SoloReversi{base, &SimpleReversiAlgolithm{}}
	}
	return &DuelReversi{base}
}

type BaseReversi struct{}

// シンプルだけど非効率なアルゴリズムになっている気がする
func (r *BaseReversi) PutableFieldCells(f domain.Field, playerStone domain.Stone) []domain.PutableFieldCell {
	field := f.Value
	opponentStone := domain.SwitchStone(playerStone)

	// 相手の駒が置かれている FieldCell を取得
	opponentFieldCellList := make([]domain.FieldCell, 0, 64-1)
	for _, row := range field {
		filtered := slices.Filter(row, func(c domain.FieldCell) bool {
			return c.Stone == opponentStone
		})
		opponentFieldCellList = append(opponentFieldCellList, filtered...)
	}

	// 左、左上、上、右上、右、右下、下、左下 の順にチェックする
	dx := []int{-1, -1, 0, 1, 1, 1, 0, -1}
	dy := []int{0, 1, 1, 1, 0, -1, -1, -1}

	// 相手の駒周囲の emptyFieldPos を取得
	emptyFieldCell := make([]domain.FieldCell, 0, 64-1)
	checked := make(map[string]domain.FieldCell)
	for _, cell := range opponentFieldCellList {
		for i := 0; i < 8; i++ {
			x, y := cell.Pos.X+dx[i], cell.Pos.Y+dy[i]
			if x < 0 || x >= 8 || y < 0 || y >= 8 {
				continue
			}

			checking, key := field[x][y], fmt.Sprintf("%d%d", x, y)
			if _, ok := checked[key]; !ok && checking.Stone == domain.EmptyStone {
				emptyFieldCell = append(emptyFieldCell, field[x][y])
				checked[key] = checking
			}
		}
	}

	// putable かバリデーション
	type putable struct {
		found bool
		arr   []domain.FieldCell
	}
	var findPutable func(x, y, dx, dy int, arr []domain.FieldCell) putable
	findPutable = func(x, y, dx, dy int, arr []domain.FieldCell) putable {
		nx, ny := x+dx, y+dy
		if nx < 0 || nx >= 8 || ny < 0 || ny >= 8 {
			return putable{found: false, arr: arr}
		}

		stone := field[nx][ny].Stone
		if stone == domain.EmptyStone {
			return putable{found: false, arr: arr}
		}

		if stone == playerStone {
			return putable{found: len(arr) != 0, arr: arr} // 隣り合うコマをカウントしないために length を見る
		}

		arr = append(arr, field[nx][ny])
		return findPutable(nx, ny, dx, dy, arr)
	}

	pfcells := make([]domain.PutableFieldCell, 0, 64-1)
	for _, cell := range emptyFieldCell {
		var pfcell domain.PutableFieldCell

		for i := 0; i < 8; i++ {
			arr := make([]domain.FieldCell, 0, 8-1)
			if putable := findPutable(cell.Pos.X, cell.Pos.Y, dx[i], dy[i], arr); putable.found {
				pfcell.FieldCell = cell
				pfcell.Stone = playerStone
				pfcell.ReversibleCells = append(
					putable.arr,
					pfcell.ReversibleCells...,
				)
			}
		}

		if len(pfcell.ReversibleCells) != 0 {
			// コマ自身も ReversibleCells に追加する
			pfcell.ReversibleCells = append(pfcell.ReversibleCells, cell)
			pfcells = append(pfcells, pfcell)
		}
	}

	return pfcells
}

func (r *BaseReversi) GetScore(f domain.Field) domain.Score {
	field := f.Value

	black, white := 0, 0
	for _, row := range field {
		for _, cell := range row {
			switch cell.Stone {
			case domain.BlackStone:
				black++
			case domain.WhiteStone:
				white++
			}
		}
	}

	var winnerStone domain.Stone
	if black > white {
		winnerStone = domain.BlackStone
	} else if white > black {
		winnerStone = domain.WhiteStone
	} else {
		winnerStone = domain.EmptyStone
	}

	return domain.Score{Black: black, White: white, WinnerStone: winnerStone}
}

func (r *BaseReversi) IsFinished(field domain.Field) bool {
	b := r.PutableFieldCells(field, domain.BlackStone)
	w := r.PutableFieldCells(field, domain.WhiteStone)
	return len(b) == 0 && len(w) == 0
}
