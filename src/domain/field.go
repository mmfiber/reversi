package domain

type Field struct {
	Value [][]FieldCell
}

func NewField() *Field {
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

	return &Field{field}
}

type FieldCell struct {
	Stone Stone
	Pos   FieldPos
}

type PutableFieldCell struct {
	FieldCell
	PutableStone    Stone
	ReversibleCells []FieldCell
}

type FieldPos struct {
	X int
	Y int
}
