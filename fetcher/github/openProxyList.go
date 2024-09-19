package github

import (
	"IpProxyPool/middleware/database"
)

func OpenProxyList() []*database.IP {
	list := make([]*database.IP, 0)
	name := "OpenProxyList"

	// http
	http := fetch(NewProxyWeb(name, "https://api.openproxylist.xyz/http.txt"))
	list = append(list, http...)
	// socks5
	socks5 := fetch(NewProxyWeb(name, "https://api.openproxylist.xyz/socks5.txt"))
	list = append(list, socks5...)
	//	 socks4
	socks4 := fetch(NewProxyWeb(name, "https://api.openproxylist.xyz/socks4.txt"))
	list = append(list, socks4...)
	return list
}
