package github

import (
	"IpProxyPool/middleware/database"
)

func Zaeem20() []*database.IP {
	list := make([]*database.IP, 0)
	name := "Zaeem20"
	// http
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/http.txt"))...)
	// https
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/https.txt"))...)
	// socks4
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/socks4.txt"))...)
	// socks5 := fetch("Zaeem20", setProxyWeb("https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/socks5.txt"))
	// list = append(list, socks5...)

	return list
}
