package gui

import (
	"reversi/src/domain"
	"sync"
)

type GuiState struct {
	sync.RWMutex

	Field              domain.Field
	CurrentPlayerStone domain.Stone
}

func NewGuiState() *GuiState {
	return &GuiState{
		Field:              domain.NewField(),
		CurrentPlayerStone: domain.BlackStone,
	}
}

func (g *GuiState) updateFieldCell(cell domain.FieldCell) {
	g.Field.Value[cell.Pos.X][cell.Pos.Y] = cell
}

func (g *GuiState) updateFieldCells(cells []domain.FieldCell) {
	g.Lock()
	defer g.Unlock()

	for _, cell := range cells {
		g.updateFieldCell(cell)
	}
}

func (g *GuiState) reverseFieldCells(putableFieldCell domain.PutableFieldCell) {
	g.Lock()
	defer g.Unlock()

	for _, cell := range putableFieldCell.ReversibleCells {
		cell.Stone = putableFieldCell.Stone
		g.updateFieldCell(cell)
	}
}

func (g *GuiState) switchStone() {
	g.CurrentPlayerStone = domain.SwitchStone(g.CurrentPlayerStone)
}
