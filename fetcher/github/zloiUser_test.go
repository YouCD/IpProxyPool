package github

import (
	"context"
	"fmt"
	"testing"
)

func TestHideIPMe(t *testing.T) {
	for _, ip := range ZloiUser(context.Background()) {
		fmt.Println(ip)
	}
}
