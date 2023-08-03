package gui

import (
	"reversi/src/domain"
	"reversi/src/utility/strconverter"
	"strconv"

	"github.com/rivo/tview"
)

type FieldView struct {
	*tview.Table
	GuiView
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
