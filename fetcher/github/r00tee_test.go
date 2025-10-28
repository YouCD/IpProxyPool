package github

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/youcd/toolkit/log"
)

func TestR00tee(t *testing.T) {
	now := time.Now()
	a := R00tee(context.Background())
	for _, ip := range a {
		log.Infof("%#v", ip)
	}
	fmt.Println(len(a))
	log.Infof("Duration: %s", time.Since(now))
}
