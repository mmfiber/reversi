package gui

import (
	"fmt"
	"reversi/src/domain"

	"github.com/rivo/tview"
)

type GameScoreView struct {
	*tview.TextView
	GuiView
}

func newGameScoreView() *GameScoreView {
	g := &GameScoreView{TextView: tview.NewTextView()}
	g.SetDisabled(true)
	return g
}

func (gs *GameScoreView) update(g *Gui) {
	var result string

	score := g.reversi.GetScore()
	switch score.WinnerStone {
	case domain.BlackStone:
		if g.reversi.IsSoloPlay() {
			result = "Win player"
		} else {
			result = "Win player 1"
		}
	case domain.WhiteStone:
		if g.reversi.IsSoloPlay() {
			result = "Win computer"
		} else {
			result = "Win player 2"
		}
	default:
		result = "Even"
	}

	g.Application.QueueUpdateDraw(func() {
		gs.Clear()
		fmt.Fprintf(gs, "Result: %s\n", result)
		fmt.Fprintf(gs, "Score:\n")
		fmt.Fprintf(gs, "\tPlayer1: %d\n", score.Black)
		fmt.Fprintf(gs, "\tPlayer2: %d\n", score.White)
	})
}
