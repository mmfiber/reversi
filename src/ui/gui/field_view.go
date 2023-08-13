package gui

import (
	"reversi/src/domain"
	"reversi/src/utility/strconverter"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FieldView struct {
	*tview.Table
	GuiView

	highlightedCell *domain.FieldCell
}

func newFieldView() *FieldView {
	emptyFieldView := func() *tview.Table {
		table := tview.NewTable()
		for r := 0; r < 9; r++ {
			row := r * 2
			for c := 0; c < 9; c++ {
				col := c * 2
				char := EMPTY_STONE_UNICODE
				if r == 0 && c != 0 {
					char, _ = strconverter.IntToCapitalizedChar(r + c)
				}
				if r != 0 && c == 0 {
					char = strconv.Itoa(r + c)
				}
				cell := tview.NewTableCell(char).SetAlign(tview.AlignCenter)
				table.SetCell(row, col, cell)
				table.SetCell(row+1, col, tview.NewTableCell("\u23E4"))   // devider for row
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
	hcell := f.highlightedCell
	for ridx, row := range field.Value {
		for cidx, cell := range row {
			var newcell *tview.TableCell
			switch cell.Stone {
			case domain.BlackStone:
				newcell = tview.NewTableCell(BLACK_STONE_UNICODE)
			case domain.WhiteStone:
				newcell = tview.NewTableCell(WHITE_STONE_UNICODE)
			default:
				newcell = tview.NewTableCell(EMPTY_STONE_UNICODE)
			}

			if hcell != nil && hcell.Pos.X == ridx && hcell.Pos.Y == cidx {
				newcell.SetBackgroundColor(tcell.ColorAqua)
			}

			// 上部・左部のインデックスとボーダーを考慮して cell を update
			g.Application.QueueUpdateDraw(func() {
				f.Table.SetCell((ridx+1)*2, (cidx+1)*2, newcell)
			})
		}
	}
}
