package ip89

import (
	"fmt"
	"testing"
)

func TestIp89(t *testing.T) {
	for _, ip := range IP89() {
		fmt.Printf("%#v\n", ip)
	}
}
