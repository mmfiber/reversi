package time

import (
	"math/rand"
	"time"
)

func Wait(min, max int) {
	rand.NewSource(time.Now().UnixNano()) // 乱数のシードを現在時刻で初期化

	diff := max - min
	if diff < 0 {
		diff = 0
	}
	sec := time.Duration(min+rand.Intn(diff)) * time.Second
	time.Sleep(sec)
}
