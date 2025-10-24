package ip3366

import (
	"fmt"
	"testing"
)

func TestIp3366(t *testing.T) {
	for _, ip := range Ip3366() {
		fmt.Printf("%#v\n", ip)
	}
}
