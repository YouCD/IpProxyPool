package github

import (
	"context"
	"fmt"
	"testing"
)

func TestYemixzy(t *testing.T) {
	for _, ip := range Yemixzy(context.Background()) {
		fmt.Println(ip)
	}
}
