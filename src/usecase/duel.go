package usecase

type DuelReversiStrategy struct {
	Reversi
}

func (h *DuelReversiStrategy) onPostPutOrPass(r *Reversi) {}
