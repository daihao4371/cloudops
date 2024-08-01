package common

import (
	"math/rand"
	"time"
)

// Frn 函数返回一个小于n的随机整数
func Frn(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}
