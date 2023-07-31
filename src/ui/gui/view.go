package gui

import (
	"fmt"
	"reversi/src/domain"
	"reversi/src/utility/strconverter"
	"strconv"

	"github.com/rivo/tview"
)

type FieldView struct {
	*tview.Table
}

func newFieldView() *FieldView {
	emptyFieldView := func() *tview.Table {
		table := tview.NewTable()
		for r := 0; r < 9; r++ {
			row := r * 2
			for c := 0; c < 9; c++ {
				col := c * 2
				char := " "
				if r == 0 && c != 0 {
					char, _ = strconverter.IntToCapitalizedChar(r + c)
				}
				if r != 0 && c == 0 {
					char = strconv.Itoa(r + c)
				}
				cell := tview.NewTableCell(char).SetAlign(tview.AlignCenter)
				table.SetCell(row, col, cell)
				table.SetCell(row+1, col, tview.NewTableCell("\u2500"))   // devider for row
				table.SetCell(row, col+1, tview.NewTableCell("\u2502"))   // devider for col
				table.SetCell(row+1, col+1, tview.NewTableCell("\u253c")) // join for row and col devider
			}
		}
		return table
	}

	f := &FieldView{Table: emptyFieldView()}
	f.SetBorder(false).SetTitle("Field").SetTitleAlign(tview.AlignLeft)

	return f
}

func (f *FieldView) update(g *Gui) {
	field := g.reversi.Field()
	for ridx, row := range field.Value {
		for cidx, cell := range row {
			var newcell *tview.TableCell
			switch cell.Stone {
			case domain.BlackStone:
				newcell = tview.NewTableCell("\u26AB")
			case domain.WhiteStone:
				newcell = tview.NewTableCell("\u26AA")
			default:
				newcell = tview.NewTableCell(" ")
			}
			// 上部・左部のインデックスとボーダーを考慮して cell を update
			f.Table.SetCell((ridx+1)*2, (cidx+1)*2, newcell)
		}
	}
}

type FieldCellSelectorView struct {
	*tview.List
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
				g.Application.Stop()
			},
		)
	})
}
