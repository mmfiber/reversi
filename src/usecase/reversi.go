package usecase

import (
	"fmt"
	"reversi/src/domain"
	"reversi/src/utility/slices"
)

type ReversiStrategy interface {
	postPut(r *Reversi)
}

func getRevesiHandler(solo, duel bool) ReversiStrategy {
	if solo {
		return &SoloReversiStrategy{&SimpleReversiAlgolithm{}}
	}
	return &DuelReversiStrategy{}
}

type Reversi struct {
	strategy           ReversiStrategy
	field              *domain.Field
	currentPlayerStone domain.Stone
}

func NewReversi(solo, duel bool) *Reversi {
	return &Reversi{
		strategy:           getRevesiHandler(solo, duel),
		field:              domain.NewField(),
		currentPlayerStone: domain.BlackStone,
	}
}

func (r *Reversi) Field() *domain.Field {
	return r.field
}

func (r *Reversi) CurrentPlayerStone() domain.Stone {
	return r.currentPlayerStone
}

// シンプルだけど非効率なアルゴリズムになっている気がする
func (r *Reversi) PutableFieldCells(playerStone domain.Stone) []domain.PutableFieldCell {
	field := r.field.Value
	opponentStone := domain.SwitchStone(playerStone)

	// 相手の駒が置かれている FieldCell を取得
	opponentFieldCellList := make([]*domain.FieldCell, 0, 64-1)
	for _, row := range field {
		filtered := slices.Filter(row, func(c *domain.FieldCell) bool {
			return c.Stone == opponentStone
		})
		opponentFieldCellList = append(opponentFieldCellList, filtered...)
	}

	// 左、左上、上、右上、右、右下、下、左下 の順にチェックする
	dx := []int{-1, -1, 0, 1, 1, 1, 0, -1}
	dy := []int{0, 1, 1, 1, 0, -1, -1, -1}

	// 相手の駒周囲の emptyFieldPos を取得
	emptyFieldCell := make([]*domain.FieldCell, 0, 64-1)
	checked := make(map[string]*domain.FieldCell)
	for _, cell := range opponentFieldCellList {
		for i := 0; i < 8; i++ {
			x, y := cell.Pos.X+dx[i], cell.Pos.Y+dy[i]
			if x < 0 || x >= 8 || y < 0 || y >= 8 {
				continue
			}

			checking, key := field[x][y], fmt.Sprintf("%d%d", x, y)
			if checking.Stone == domain.EmptyStone && checked[key] == nil {
				emptyFieldCell = append(emptyFieldCell, field[x][y])
				checked[key] = checking
			}
		}
	}

	// palceable かバリデーション
	var putable func(x, y, dx, dy int, arr *[]*domain.FieldCell) bool
	putable = func(x, y, dx, dy int, arr *[]*domain.FieldCell) bool {
		nx, ny := x+dx, y+dy
		if nx < 0 || nx >= 8 || ny < 0 || ny >= 8 {
			return false
		}

		stone := field[nx][ny].Stone
		if stone == domain.EmptyStone {
			return false
		}

		if stone == playerStone {
			return len(*arr) != 0 // 隣り合うコマをカウントしないために length を見る
		}

		*arr = append(*arr, field[nx][ny])
		return putable(nx, ny, dx, dy, arr)
	}

	putableFieldCells := make([]domain.PutableFieldCell, 0, 64-1)
	for _, cell := range emptyFieldCell {
		var putableFieldCell *domain.PutableFieldCell

		for i := 0; i < 8; i++ {
			arr := make([]*domain.FieldCell, 0, 8-1)
			if putable(cell.Pos.X, cell.Pos.Y, dx[i], dy[i], &arr) {
				if putableFieldCell != nil {
					putableFieldCell.ReversiableCells = append(arr, putableFieldCell.ReversiableCells...)
					continue
				}

				putableFieldCell = &domain.PutableFieldCell{
					FieldCell:        cell,
					PutableStone:     playerStone,
					ReversiableCells: arr,
				}
			}
		}

		if putableFieldCell != nil {
			putableFieldCells = append(putableFieldCells, *putableFieldCell)
		}
	}

	return putableFieldCells
}

func (r *Reversi) Put(cell domain.PutableFieldCell) {
	field := r.field.Value
	stone := cell.PutableStone

	field[cell.Pos.X][cell.Pos.Y] = cell.FieldCell
	field[cell.Pos.X][cell.Pos.Y].Stone = stone

	for _, reversed := range cell.ReversiableCells {
		reversed.Stone = stone
		x, y := reversed.Pos.X, reversed.Pos.Y
		field[x][y] = reversed
	}

	r.currentPlayerStone = domain.SwitchStone(r.currentPlayerStone)
}

func (r *Reversi) PostPut() {
	r.strategy.postPut(r)
}

func (r *Reversi) Pass() {
	r.currentPlayerStone = domain.SwitchStone(r.currentPlayerStone)
}

func (r *Reversi) GetScore() domain.Score {
	black, white := 0, 0
	for _, row := range r.field.Value {
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

func (r *Reversi) GetFieldCell(ridx, cidx int) *domain.FieldCell {
	if ridx < 0 || ridx > 7 || cidx < 0 || cidx > 7 {
		return nil
	}

	return r.field.Value[ridx][cidx]
}

func (r *Reversi) IsFinished() bool {
	b := r.PutableFieldCells(domain.BlackStone)
	w := r.PutableFieldCells(domain.WhiteStone)
	return len(b) == 0 && len(w) == 0
}
