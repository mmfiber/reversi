package usecase

import (
	"reversi/src/domain"
)

type DuelReversiHandler struct{}

func (h *DuelReversiHandler) next(field *domain.Field, currentStone domain.Stone) domain.Stone {
	return domain.SwitchStone(currentStone)
}
