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
	n.SetRows(0).SetColumns(0)
	return n
}

func (n *NavigatorView) update(g *Gui) {
	viewSetter := func(view interface {
		tview.Primitive
		GuiView
	}) (updated bool) {
		updated = false
		if n.view == view {
			return
		}

		updated = true
		n.Clear()
		n.AddItem(view, 0, 0, 1, 1, 0, 0, true)
		n.view = view
		g.SetFocus(view)
		return
	}

	switch g.status {
	case GamePlaying:
		if viewSetter(n.fieldCellSelectorView) {
			n.fieldCellSelectorView.init(g)
		}
	case GameFinished:
		viewSetter(n.gameScoreView)
	default:
		logger.Error(fmt.Errorf("unhandleable gui status: %d", g.status))
		g.quit()
	}

	n.view.update(g)
}
