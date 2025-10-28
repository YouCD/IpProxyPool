package github

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestAnonym0usWork1221(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()
	for _, ipObj := range Anonym0usWork1221(ctx) {
		fmt.Printf("%#v\n", ipObj)
	}
}
