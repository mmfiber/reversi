package usecase

import (
	"fmt"
	"os"
	"reversi/src/domain"
	"reversi/src/log"
	"reversi/src/utility/slices"
	"sync"
)

type Reversi interface {
	Field() *domain.Field
	CurrentPlayerStone() domain.Stone
	PutableFieldCells(playerStone domain.Stone) []domain.PutableFieldCell
	Put(cell domain.PutableFieldCell)
	PostPut()
	Pass()
	PostPass()
	GetScore() domain.Score
	GetFieldCell(ridx, cidx int) domain.FieldCell
	IsFinished() bool
	IsSoloPlay() bool
}

var logger = log.NewLogger()

func NewReversi(solo, duel bool) Reversi {
	base := &BaseReversi{
		field:              domain.NewField(),
		currentPlayerStone: domain.BlackStone,
	}
	if solo {
		return &SoloReversi{base, &SimpleReversiAlgolithm{}}
	}
	return &DuelReversi{base}
}

type BaseReversi struct {
	sync.RWMutex

	field              *domain.Field
	currentPlayerStone domain.Stone
}

func (r *BaseReversi) Field() *domain.Field {
	return r.field
}

func (r *BaseReversi) CurrentPlayerStone() domain.Stone {
	return r.currentPlayerStone
}

// シンプルだけど非効率なアルゴリズムになっている気がする
func (r *BaseReversi) PutableFieldCells(playerStone domain.Stone) []domain.PutableFieldCell {
	r.RLock()
	field := r.field.Value
	opponentStone := domain.SwitchStone(playerStone)
	r.RUnlock()

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

	// palceable かバリデーション
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

	putableFieldCells := make([]domain.PutableFieldCell, 0, 64-1)
	for _, cell := range emptyFieldCell {
		var putableFieldCell domain.PutableFieldCell

		for i := 0; i < 8; i++ {
			arr := make([]domain.FieldCell, 0, 8-1)
			if putable := findPutable(cell.Pos.X, cell.Pos.Y, dx[i], dy[i], arr); putable.found {
				putableFieldCell.FieldCell = cell
				putableFieldCell.PutableStone = playerStone
				putableFieldCell.ReversibleCells = append(putable.arr, putableFieldCell.ReversibleCells...)
			}
		}

		if len(putableFieldCell.ReversibleCells) != 0 {
			putableFieldCells = append(putableFieldCells, putableFieldCell)
		}
	}

	return putableFieldCells
}

func (r *BaseReversi) Put(cell domain.PutableFieldCell) {
	r.Lock()
	defer r.Unlock()

	field := r.field.Value
	stone := cell.PutableStone

	field[cell.Pos.X][cell.Pos.Y] = cell.FieldCell
	field[cell.Pos.X][cell.Pos.Y].Stone = stone

	for _, reversed := range cell.ReversibleCells {
		reversed.Stone = stone
		x, y := reversed.Pos.X, reversed.Pos.Y
		field[x][y] = reversed
	}

	r.currentPlayerStone = domain.SwitchStone(r.currentPlayerStone)
}

func (r *BaseReversi) Pass() {
	r.Lock()
	defer r.Unlock()
	r.currentPlayerStone = domain.SwitchStone(r.currentPlayerStone)
}

func (r *BaseReversi) GetScore() domain.Score {
	r.RLock()
	field := r.field.Value
	r.RUnlock()

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

func (r *BaseReversi) GetFieldCell(ridx, cidx int) domain.FieldCell {
	r.RLock()
	defer r.RUnlock()

	if ridx < 0 || ridx > 7 || cidx < 0 || cidx > 7 {
		logger.Error(fmt.Errorf("index out of range, ridx: %d, cidx: %d", ridx, cidx))
		os.Exit(1)
	}

	return r.field.Value[ridx][cidx]
}

func (r *BaseReversi) IsFinished() bool {
	r.RLock()
	defer r.RUnlock()

	b := r.PutableFieldCells(domain.BlackStone)
	w := r.PutableFieldCells(domain.WhiteStone)
	return len(b) == 0 && len(w) == 0
}
