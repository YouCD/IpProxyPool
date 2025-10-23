package github

import (
	"IpProxyPool/middleware/database"
)

func Yemixzy() []*database.IP {
	list := make([]*database.IP, 0)
	name := "yemixzy"
	// http
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/http.txt"))...)
	// socks5
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/socks5.txt"))...)
	//	 socks4
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/socks4.txt"))...)
	//	 unchecked
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/unchecked.txt"))...)
	return list
}
