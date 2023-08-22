package usecase

type DuelReversi struct {
	*BaseReversi
}

func (d *DuelReversi) PostPut() {}

func (d *DuelReversi) PostPass() {}

func (d *DuelReversi) IsSoloPlay() bool {
	return false
}
