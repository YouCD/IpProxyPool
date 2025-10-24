package ip89

import (
	"fmt"
	"testing"
)

func TestIp89(t *testing.T) {
	for _, ip := range Ip89() {
		fmt.Printf("%#v\n", ip)
	}
}
