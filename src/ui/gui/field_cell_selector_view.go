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

const (
	PASS_TEXT = "pass"
	QUIT_TEXT = "quit"
)

type FieldCellSelectorView struct {
	*tview.Grid
	GuiView

	textView *FieldCellTextView
	listView *FieldCellListView
}

func newFieldCellSelectorView() *FieldCellSelectorView {
	f := &FieldCellSelectorView{
		Grid:     tview.NewGrid(),
		textView: newFieldCellTextView(),
		listView: newFieldCellListView(),
	}

	f.SetRows(5, 0).SetColumns(0)
	f.AddItem(f.textView, 0, 0, 1, 1, 0, 0, false)
	f.AddItem(f.listView, 1, 0, 1, 1, 0, 0, true)
	return f
}

func (f *FieldCellSelectorView) init(g *Gui) {
	f.listView.enableFieldHighlight(g)
	g.SetFocus(f.listView)
}

func (f *FieldCellSelectorView) update(g *Gui) {
	f.textView.update(g)
	f.listView.update(g)
}

type FieldCellTextView struct {
	*tview.TextView
	GuiView
}

func newFieldCellTextView() *FieldCellTextView {
	f := &FieldCellTextView{TextView: tview.NewTextView()}
	f.SetDisabled(true)
	return f
}

func (f *FieldCellTextView) update(g *Gui) {
	f.Clear()

	var player, stoneUnicode string
	switch stone := g.reversi.CurrentPlayerStone(); stone {
	case domain.BlackStone:
		player = "player 1"
		stoneUnicode = BLACK_STONE_UNICODE
	case domain.WhiteStone:
		player = "player 2"
		stoneUnicode = WHITE_STONE_UNICODE
	default:
		logger.Error(fmt.Errorf("unhandleable stone: %d", stone))
		g.quit()
	}

	fmt.Fprintf(f, "Current Player: %s\n", player)
	fmt.Fprintf(f, "Player Stone: %s\n", stoneUnicode)
}

type FieldCellListView struct {
	*tview.List
	GuiView
}

func newFieldCellListView() *FieldCellListView {
	f := &FieldCellListView{List: tview.NewList()}
	f.SetTitle(" Select where to put your stone ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	return f
}

func (f *FieldCellListView) posToIndex(pos domain.FieldPos) (string, string) {
	row := fmt.Sprint(pos.X + 1)
	col, _ := strconverter.IntToCapitalizedChar(pos.Y + 1)
	return row, col
}

func (f *FieldCellListView) indexToPos(row, col string) domain.FieldPos {
	ridx, err := strconv.Atoi(row)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	cidx := strconverter.CapitalizedCharToInt(col)
	return domain.FieldPos{X: ridx - 1, Y: cidx - 1}
}

func (f *FieldCellListView) enableFieldHighlight(g *Gui) {
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

func (f *FieldCellListView) update(g *Gui) {
	if g.reversi.IsFinished() {
		g.gameFinished()
		return
	}

	f.Clear()
	g.Application.QueueUpdateDraw(func() {
		playerStone := g.reversi.CurrentPlayerStone()
		putableFieldCells := g.reversi.PutableFieldCells(playerStone)
		if len(putableFieldCells) != 0 {
			for idx, cell := range putableFieldCells {
				escapedCell := cell
				row, col := f.posToIndex(cell.Pos)
				char, _ := strconverter.IntToChar(idx + 1)
				f.AddItem(
					fmt.Sprintf("%s%s", row, col),
					fmt.Sprintf("row %s and col %s", row, col),
					strconverter.CharToRune(char),
					func() {
						g.reversi.Put(escapedCell)
						g.updateFieldView()
						g.reversi.PostPut()
						g.updateView()
					},
				)
			}
			// 先頭の選択肢をハイライト
			g.highlightFieldCell(putableFieldCells[0].FieldCell)
		} else {
			f.AddItem(
				"pass",
				"pass",
				'p',
				func() {
					g.reversi.Pass()
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
