package github

import (
	"fmt"
	"testing"
)

func TestFreshProxyList(t *testing.T) {
	for _, ip := range Vakhov() {
		fmt.Println(ip)
	}
}
