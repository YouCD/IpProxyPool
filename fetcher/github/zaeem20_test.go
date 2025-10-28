package github

import (
	"context"
	"fmt"
	"testing"
)

func TestZaeem20(t *testing.T) {
	for _, ip := range Zaeem20(context.Background()) {
		fmt.Printf("%#v\n", ip)
	}
}
