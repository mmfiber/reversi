package usecase

import (
	"fmt"
	"reversi/src/domain"
	"reversi/src/utility/slices"
)

type ReversiHandler interface {
	next(field *domain.Field, currentStone domain.Stone) domain.Stone
}

func getRevesiHandler() ReversiHandler {
	return &DuelReversiHandler{}
}

type Reversi struct {
	handler            ReversiHandler
	field              *domain.Field
	currentPlayerStone domain.Stone
}

func NewReversi() *Reversi {
	return &Reversi{
		handler:            getRevesiHandler(),
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
func (r *Reversi) PlaceableFieldCells(playerStone domain.Stone) []domain.PlaceableFieldCell {
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
	var placeable func(x, y, dx, dy int, arr *[]*domain.FieldCell) bool
	placeable = func(x, y, dx, dy int, arr *[]*domain.FieldCell) bool {
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
		return placeable(nx, ny, dx, dy, arr)
	}

	placeableFieldCells := make([]domain.PlaceableFieldCell, 0, 64-1)
	for _, cell := range emptyFieldCell {
		var placeableFieldCell *domain.PlaceableFieldCell

		for i := 0; i < 8; i++ {
			arr := make([]*domain.FieldCell, 0, 8-1)
			if placeable(cell.Pos.X, cell.Pos.Y, dx[i], dy[i], &arr) {
				if placeableFieldCell != nil {
					placeableFieldCell.ReversiableCells = append(arr, placeableFieldCell.ReversiableCells...)
					continue
				}

				placeableFieldCell = &domain.PlaceableFieldCell{
					FieldCell:        cell,
					PlaceableStone:   playerStone,
					ReversiableCells: arr,
				}
			}
		}

		if placeableFieldCell != nil {
			placeableFieldCells = append(placeableFieldCells, *placeableFieldCell)
		}
	}

	return placeableFieldCells
}

func (r *Reversi) Placed(cell domain.PlaceableFieldCell) {
	field := r.field.Value
	stone := cell.PlaceableStone

	placedFieldCell := cell.FieldCell
	placedFieldCell.Stone = stone
	field[cell.Pos.X][cell.Pos.Y] = placedFieldCell

	for _, reversed := range cell.ReversiableCells {
		reversed.Stone = stone
		x, y := reversed.Pos.X, reversed.Pos.Y
		field[x][y] = reversed
	}

	r.currentPlayerStone = r.handler.next(r.field, stone)
}

func (r *Reversi) Pass(currentStone domain.Stone) {
	r.currentPlayerStone = r.handler.next(r.field, currentStone)
}
