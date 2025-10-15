package database

import (
	"IpProxyPool/middleware/config"
	"fmt"
	"github.com/youcd/toolkit/log"
	"testing"
)

func TestDeleteByIP(t *testing.T) {
	host := GetIPByProxyHost("8.210.34.11")
	fmt.Printf("%v", host)
}
