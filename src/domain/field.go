package domain

import (
	"fmt"
	"os"
)

type Field struct {
	Value [][]FieldCell
}

func NewField() Field {
	field := make([][]FieldCell, 8)
	for x := range field {
		field[x] = make([]FieldCell, 8)
		for y := range field[x] {
			field[x][y] = FieldCell{
				Stone: EmptyStone,
				Pos:   FieldPos{x, y},
			}
		}
	}

	field[3][3].Stone = WhiteStone
	field[3][4].Stone = BlackStone
	field[4][3].Stone = BlackStone
	field[4][4].Stone = WhiteStone

	return Field{field}
}

func (f *Field) GetCell(ridx, cidx int) FieldCell {
	if ridx < 0 || ridx > 7 || cidx < 0 || cidx > 7 {
		logger.Error(fmt.Errorf("index out of range, ridx: %d, cidx: %d", ridx, cidx))
		os.Exit(1)
	}

	return f.Value[ridx][cidx]
}

type FieldCell struct {
	Stone Stone
	Pos   FieldPos
}

type PutableFieldCell struct {
	FieldCell
	Stone           Stone
	ReversibleCells []FieldCell
}

type FieldPos struct {
	X int
	Y int
}
