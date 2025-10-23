package github

import (
	"IpProxyPool/middleware/database"
)

func Vakhov() []*database.IP {
	list := make([]*database.IP, 0)
	name := "Vakhov"

	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/https.txt"))...)

	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/socks5.txt"))...)

	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/socks4.txt"))...)

	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/http.txt"))...)
	return list
}
