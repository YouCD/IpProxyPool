package proxylistplus

import (
	"fmt"
	"testing"
)

func TestProxyListPlus(t *testing.T) {
	for _, ip := range ProxyListPlus() {
		fmt.Printf("%#v\n", ip)
	}
}
