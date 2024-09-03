package github

import "IpProxyPool/middleware/database"

func TheSpeedX() []*database.IP {
	list := make([]*database.IP, 0)

	httpStr := setProxyWeb("https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt")
	list = append(list, fetch("ProxyList", httpStr)...)

	socks4Str := setProxyWeb("https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks4.txt")
	list = append(list, fetch("ProxyList", socks4Str)...)

	socks5Str := setProxyWeb("https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks5.txt")
	list = append(list, fetch("ProxyList", socks5Str)...)
	return list
}
