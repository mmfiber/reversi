package gui

import (
	"fmt"
	"os"
	"reversi/src/domain"
	"reversi/src/utility/strconverter"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FieldCellSelectorView struct {
	*tview.List
	GuiView
}

const (
	PASS_TEXT = "pass"
	QUIT_TEXT = "quit"
)

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

func (f *FieldCellSelectorView) indexToPos(row, col string) domain.FieldPos {
	ridx, err := strconv.Atoi(row)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	cidx := strconverter.CapitalizedCharToInt(col)
	return domain.FieldPos{X: ridx - 1, Y: cidx - 1}
}

func (f *FieldCellSelectorView) enableFieldHighlight(g *Gui) {
	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		var next int
		count := f.GetItemCount()
		current := f.GetCurrentItem()
		switch event.Key() {
		case tcell.KeyUp:
			if current >= 1 {
				next = current - 1
			}
		case tcell.KeyDown:
			if current < count-1 {
				next = current + 1
			}
		case tcell.KeyLeft:
			if current >= 1 {
				next = current - 1
			}
		case tcell.KeyRight:
			if current < count-1 {
				next = current + 1
			}
		default:
			return event
		}

		mainText, _ := f.GetItemText(next)
		if mainText == PASS_TEXT || mainText == QUIT_TEXT {
			return event
		}
		if len(mainText) != 2 {
			return event
		}

		pos := f.indexToPos(mainText[:1], mainText[1:])
		cell := g.reversi.GetFieldCell(pos.X, pos.Y)
		g.highlightFieldCell(cell)

		return event
	})
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
			// 先頭の選択肢をハイライト
			g.highlightFieldCell(placeableFieldCells[0].FieldCell)
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
