package github

import (
	"context"
	"testing"

	"github.com/youcd/toolkit/log"
)

func init() {

}
func TestFreeProxyList(t *testing.T) {
	for _, ip := range FreeProxyList(context.Background()) {
		log.Infof("%#v", ip)
	}
}
