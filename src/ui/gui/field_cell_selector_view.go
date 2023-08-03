package gui

import (
	"fmt"
	"reversi/src/domain"
	"reversi/src/utility/strconverter"

	"github.com/rivo/tview"
)

type FieldCellSelectorView struct {
	*tview.List
	GuiView
}

func newFieldCellSelectorView() *FieldCellSelectorView {
	f := &FieldCellSelectorView{List: tview.NewList()}
	f.SetTitle(" Select where to put your stone ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	return f
}

func (f *FieldCellSelectorView) posToIndex(pos domain.FieldPos) (string, string) {
	row := fmt.Sprint(pos.X + 1)
	col, _ := strconverter.IntToCapitalizedChar(pos.Y + 1)
	return row, col
}

func (f *FieldCellSelectorView) update(g *Gui) {
	if g.reversi.IsFinished() {
		g.gameFinished()
		return
	}

	f.Clear()
	g.Application.QueueUpdateDraw(func() {
		playerStone := g.reversi.CurrentPlayerStone()
		placeableFieldCells := g.reversi.PlaceableFieldCells(playerStone)
		if len(placeableFieldCells) != 0 {
			for idx, cell := range placeableFieldCells {
				escapedCell := cell
				row, col := f.posToIndex(cell.Pos)
				char, _ := strconverter.IntToChar(idx + 1)
				f.AddItem(
					fmt.Sprintf("%s%s", row, col),
					fmt.Sprintf("row %s and col %s", row, col),
					strconverter.CharToRune(char),
					func() {
						g.reversi.Placed(escapedCell)
						g.updateView()
					},
				)
			}
		} else {
			f.AddItem(
				"pass",
				"pass",
				'p',
				func() {
					g.reversi.Pass(playerStone)
					g.updateView()
				},
			)
		}
		f.AddItem(
			"quit",
			"finish and quit the game",
			'q',
			func() {
				g.quit()
			},
		)
	})
}
