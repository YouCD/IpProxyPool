package github

import (
	"IpProxyPool/middleware/database"
)

func Zaeem20() []*database.IP {
	list := make([]*database.IP, 0)
	name := "Zaeem20"
	// http
	http := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/http.txt"))
	list = append(list, http...)
	// https
	https := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/https.txt"))
	list = append(list, https...)
	// socks4
	socks4 := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/socks4.txt"))
	list = append(list, socks4...)
	// socks5 := fetch("Zaeem20", setProxyWeb("https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/socks5.txt"))
	// list = append(list, socks5...)

	return list
}
