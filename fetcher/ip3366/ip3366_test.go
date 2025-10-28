package ip3366

import (
	"context"
	"fmt"
	"testing"
)

func TestIp3366(t *testing.T) {
	for _, ip := range Ip3366(context.Background()) {
		fmt.Printf("%#v\n", ip)
	}
}
