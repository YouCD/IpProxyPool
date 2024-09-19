package github

import (
	"IpProxyPool/middleware/database"
)

func TheSpeedX() []*database.IP {
	list := make([]*database.IP, 0)

	name := "TheSpeedX"

	httpStr := NewProxyWeb(name, "https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt")
	list = append(list, fetch(httpStr)...)

	socks4Str := NewProxyWeb(name, "https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks4.txt")
	list = append(list, fetch(socks4Str)...)

	socks5Str := NewProxyWeb(name, "https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks5.txt")
	list = append(list, fetch(socks5Str)...)
	return list
}
