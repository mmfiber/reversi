package gui

import (
	"os"
	"reversi/src/domain"
	"reversi/src/log"
	"reversi/src/usecase"

	"github.com/rivo/tview"
)

var logger = log.NewLogger()

type Gui struct {
	*tview.Application
	reversi       *usecase.Reversi
	fieldView     *FieldView
	navigatorView *NavigatorView
	status        GuiStatus
}

type GuiView interface {
	update(*Gui)
}

type GuiStatus = int

const (
	GamePlaying GuiStatus = iota
	GameFinished
	GameQuit
)

func New(solo, duel bool) *Gui {
	return &Gui{
		Application:   tview.NewApplication(),
		reversi:       usecase.NewReversi(solo, duel),
		fieldView:     newFieldView(),
		navigatorView: newNavigatorView(),
		status:        GamePlaying,
	}
}

func (g *Gui) Run() {
	divider := tview.NewBox()
	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(-1, 10, -2).
		SetOffset(0, 1).
		AddItem(g.navigatorView, 0, 0, 1, 1, 0, 0, true).
		AddItem(divider, 0, 0, 1, 1, 0, 0, true).
		AddItem(g.fieldView, 0, 2, 1, 1, 0, 0, false)

	g.updateView()

	if err := g.SetRoot(grid, true).SetFocus(g.navigatorView).Run(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func (g *Gui) updateView() {
	views := [...]GuiView{
		g.fieldView,
		g.navigatorView,
	}
	for _, view := range views {
		go view.update(g)
	}
}

func (g *Gui) updateFieldView() {
	go g.fieldView.update(g)
}

func (g *Gui) highlightFieldCell(cell *domain.FieldCell) {
	g.fieldView.highlightedCell = cell
	g.updateFieldView()
}

func (g *Gui) gameFinished() {
	g.status = GameFinished
	g.updateView()
}

func (g *Gui) quit() {
	g.status = GameQuit
	g.Application.Stop()
}
