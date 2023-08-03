package gui

import (
	"fmt"

	"github.com/rivo/tview"
)

type NavigatorView struct {
	*tview.Grid
	GuiView

	view                  GuiView
	gameScoreView         *GameScoreView
	fieldCellSelectorView *FieldCellSelectorView
}

func newNavigatorView() *NavigatorView {
	n := &NavigatorView{
		Grid:                  tview.NewGrid(),
		gameScoreView:         newGameScoreView(),
		fieldCellSelectorView: newFieldCellSelectorView(),
	}
	n.SetRows(1).SetColumns(1)
	return n
}

func (n *NavigatorView) update(g *Gui) {
	viewSetter := func(view interface {
		tview.Primitive
		GuiView
	}) {
		if n.view == view {
			return
		}
		n.Clear()
		n.AddItem(view, 1, 1, 1, 1, 0, 0, true)
		n.view = view
		g.SetFocus(view)
	}

	switch g.status {
	case GamePlaying:
		viewSetter(n.fieldCellSelectorView)
	case GameFinished:
		viewSetter(n.gameScoreView)
	default:
		logger.Error(fmt.Errorf("unhandleable gui status: %d", g.status))
		g.quit()
	}

	n.view.update(g)
}
