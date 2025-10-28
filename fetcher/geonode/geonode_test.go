package geonode

import (
	"context"
	"fmt"
	"testing"
)

func TestGeonode(t *testing.T) {
	for _, ip := range Geonode(context.Background()) {
		fmt.Printf("%#v\n", ip)
	}
}
