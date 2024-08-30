package storage

import (
	"IpProxyPool/middleware/config"
	"IpProxyPool/middleware/database"
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
	database.InitDB(&setting.Database)
}
func TestCheckProxyDB(t *testing.T) {

	CheckProxyDB()
}

func TestCheckIp(t *testing.T) {
	ip1 := &database.IP{
		ProxyId:       80,
		ProxyHost:     "47.115.219.60",
		ProxyPort:     7890,
		ProxyType:     "HTTPS",
		ProxyLocation: "SSL高匿_中国阿里云",
		ProxySpeed:    1046,
		ProxySource:   "http://www.ip3366.net",
		CreateTime:    "",
		UpdateTime:    "",
	}
	ip2 := &database.IP{
		ProxyId:       4842,
		ProxyHost:     "81.69.33.240",
		ProxyPort:     7890,
		ProxyType:     "HTTPS",
		ProxyLocation: "SSL高匿_上海市腾讯云",
		ProxySpeed:    1046,
		ProxySource:   "http://www.ip3366.net",
		CreateTime:    "",
		UpdateTime:    "",
	}
	ip3 := &database.IP{
		ProxyId:       4842,
		ProxyHost:     "222.190.173.176",
		ProxyPort:     8089,
		ProxyType:     "https",
		ProxyLocation: "SSL高匿_上海市腾讯云",
		ProxySpeed:    1046,
		ProxySource:   "http://www.ip3366.net",
		CreateTime:    "",
		UpdateTime:    "",
	}

	fmt.Println(ip2)
	fmt.Println(ip1)
	fmt.Println(ip3)
	fmt.Println(CheckIp(ip3))

}
