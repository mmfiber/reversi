package usecase

import (
	"reversi/src/domain"
)

type DuelReversiHandler struct{}

func (h *DuelReversiHandler) postPut(r *Reversi) domain.Stone {
	return domain.SwitchStone(r.currentPlayerStone)
}
