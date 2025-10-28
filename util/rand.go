//nolint:revive
package util

import (
	"math/rand"
	"time"
)

// RandInt get the random number in [min, max]
func RandInt(miniValue, maxValue int) int {
	if miniValue >= maxValue || maxValue == 0 {
		return maxValue
	}
	//nolint:gosec
	rand.New(rand.NewSource(time.Now().Local().UnixNano()))
	//nolint:gosec
	num := rand.Intn(maxValue-miniValue) + miniValue
	return num
}
