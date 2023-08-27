package usecase

import "reversi/src/domain"

type DuelReversi struct {
	*BaseReversi
}

func (d *DuelReversi) Put(field domain.Field, playerStone domain.Stone) domain.PutableFieldCell {
	// 2人プレイ中にオセロのコマを置くロジックを定義したくなった場合にここに書く
	return domain.PutableFieldCell{}
}

func (d *DuelReversi) IsSoloPlay() bool {
	return false
}
