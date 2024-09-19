package github

import (
	"fmt"
	"testing"
)

func TestYemixzy(t *testing.T) {
	for _, ip := range Yemixzy() {
		fmt.Println(ip)
	}
}
