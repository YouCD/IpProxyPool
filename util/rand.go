package util

import (
	"math/rand"
	"time"
)

// RandInt get the random number in [min, max]
//
//nolint:gosec,predeclared
func RandInt(mini, max int) int {
	if mini >= max || max == 0 {
		return max
	}
	rand.New(rand.NewSource(time.Now().Local().UnixNano()))
	num := rand.Intn(max-mini) + mini
	return num
}
