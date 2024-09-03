package github

import (
	"fmt"
	"testing"
)

func TestFreshProxyList(t *testing.T) {
	for _, ip := range FreshProxyList() {
		fmt.Println(ip)
	}
}
