package common

import (
	"math/rand"
	"time"
)

func Frn(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}
