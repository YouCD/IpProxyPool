package github

import (
	"fmt"
	"github.com/youcd/toolkit/log"
	"testing"
)

func init() {
	log.Init(true)
}
func TestFreeProxyList(t *testing.T) {
	for _, ip := range FreeProxyList() {
		fmt.Println(ip)
	}
}
