package github

import (
	"fmt"
	"testing"
)

func TestZaeem20(t *testing.T) {
	for _, ip := range Zaeem20() {
		fmt.Printf("%#v\n", ip)
	}
}
