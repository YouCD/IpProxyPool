package github

import (
	"fmt"
	"testing"

	"github.com/youcd/toolkit/log"
)

func init() {
	log.Init(nil)
}
func TestFreeProxyList(t *testing.T) {
	for _, ip := range FreeProxyList() {
		fmt.Println(ip)
	}
}
