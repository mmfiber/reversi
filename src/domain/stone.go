package domain

type Stone = int

const (
	BlackStone Stone = iota - 1
	EmptyStone
	WhiteStone
)

func SwitchStone(current Stone) Stone {
	return current * -1
}
