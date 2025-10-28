package github

import (
	"context"
	"fmt"
	"testing"
)

func TestFreshProxyList(t *testing.T) {
	for _, ip := range Vakhov(context.Background()) {
		fmt.Println(ip)
	}
}
