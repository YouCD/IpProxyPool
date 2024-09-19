package github

import (
	"IpProxyPool/middleware/database"
)

func Yemixzy() []*database.IP {
	list := make([]*database.IP, 0)
	name := "yemixzy"
	// http
	http := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/http.txt"))
	list = append(list, http...)
	// socks5
	socks5 := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/socks5.txt"))
	list = append(list, socks5...)
	//	 socks4
	socks4 := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/socks4.txt"))
	list = append(list, socks4...)
	//	 unchecked
	unchecked := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/unchecked.txt"))
	list = append(list, unchecked...)
	return list
}
