package github

import (
	"fmt"
	"testing"
)

func TestHideIPMe(t *testing.T) {
	for _, ip := range HideIPMe() {
		fmt.Println(ip)
	}
}
