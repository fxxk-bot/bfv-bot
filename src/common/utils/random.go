package utils

import (
	"golang.org/x/exp/rand"
	"time"
)

func RandomInt(max int) int {
	rand.Seed(uint64(time.Now().UnixNano()))
	return rand.Intn(max)
}
