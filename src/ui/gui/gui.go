package gui

import (
	"os"
	"reversi/src/log"
	"reversi/src/usecase"

	"github.com/rivo/tview"
)

var logger = log.NewLogger()

type Gui struct {
	*tview.Application
	fieldView             *FieldView
	fieldCellSelectorView *FieldCellSelectorView
	reversi               *usecase.Reversi
}

func New() *Gui {
	return &Gui{
		Application:           tview.NewApplication(),
		fieldView:             newFieldView(),
		fieldCellSelectorView: newFieldCellSelectorView(),
		reversi:               usecase.NewReversi(),
	}
}

func (g *Gui) Run() {
	divider := tview.NewBox()
	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(-1, 10, -2).
		SetOffset(0, 1).
		AddItem(g.fieldCellSelectorView, 0, 0, 1, 1, 0, 0, true).
		AddItem(divider, 0, 0, 1, 1, 0, 0, true).
		AddItem(g.fieldView, 0, 2, 1, 1, 0, 0, false)

	g.updateView()

	if err := g.SetRoot(grid, true).SetFocus(g.fieldCellSelectorView).Run(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func (g *Gui) updateView() {
	go g.fieldView.update(g)
	go g.fieldCellSelectorView.update(g)
}
