package database

import (
	"IpProxyPool/middleware/config"
	"fmt"
	"github.com/youcd/toolkit/log"
	"testing"
)

func init() {
	config.ConfigFile = "/home/ycd/self_data/source_code/IpProxyPool/conf/config.yaml"
	config.InitConfig()
	setting := config.ServerSetting
	log.Init(true)
	log.SetLogLevel(setting.Log.Level)
	InitDB(&setting.Database)
}
func TestDeleteByIP(t *testing.T) {
	host := GetIPByProxyHost("8.210.34.11")
	fmt.Printf("%v", host)
}
