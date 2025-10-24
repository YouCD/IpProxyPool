package geonode

import (
	"fmt"
	"testing"
)

func TestGeonode(t *testing.T) {
	for _, ip := range Geonode() {
		fmt.Printf("%#v\n", ip)
	}
}
