package github

import (
	"IpProxyPool/middleware/config"
	"context"
	"fmt"
	"testing"

	"github.com/youcd/toolkit/log"
)

func init() {
	config.ConfigFile = "/home/ycd/self_data/source_code/IpProxyPool/conf/config.yaml"
	config.InitConfig()
	setting := config.ServerSetting
	log.Init(nil)
	log.SetLogLevel(setting.Log.Level)
}
func TestOpenProxyList(t *testing.T) {
	for _, ipObj := range OpenProxyList(context.Background()) {
		fmt.Printf("%#v\n", ipObj)
	}
}
