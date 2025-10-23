package github

import (
	"IpProxyPool/middleware/database"
)

func TheSpeedX() []*database.IP {
	list := make([]*database.IP, 0)

	name := "TheSpeedX"

	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt"))...)

	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks4.txt"))...)

	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks5.txt"))...)
	return list
}
