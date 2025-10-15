package github

import (
	"IpProxyPool/middleware/database"
)

func OpenProxyList() []*database.IP {
	list := make([]*database.IP, 0)
	name := "OpenProxyList"
	// https
	list = append(list, fetch(NewProxyWeb(name, "https://api.openproxylist.xyz/https.txt"))...)
	// http
	list = append(list, fetch(NewProxyWeb(name, "https://api.openproxylist.xyz/http.txt"))...)
	// socks5
	list = append(list, fetch(NewProxyWeb(name, "https://api.openproxylist.xyz/socks5.txt"))...)
	//	 socks4
	list = append(list, fetch(NewProxyWeb(name, "https://api.openproxylist.xyz/socks4.txt"))...)
	return list
}
