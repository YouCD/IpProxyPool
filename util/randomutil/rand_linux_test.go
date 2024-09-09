package randomutil

import (
	"fmt"
	"testing"
)

func TestRandInt(t *testing.T) {
	fmt.Println(RandInt(1, 30000))
}
