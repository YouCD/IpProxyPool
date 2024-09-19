package github

import (
	"IpProxyPool/middleware/database"
)

func Vakhov() []*database.IP {
	list := make([]*database.IP, 0)
	name := "Vakhov"

	httpsURL := NewProxyWeb(name, "https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/https.txt")
	list = append(list, fetch(httpsURL)...)

	socks5Url := NewProxyWeb(name, "https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/socks5.txt")
	list = append(list, fetch(socks5Url)...)

	socks4Url := NewProxyWeb(name, "https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/socks4.txt")
	list = append(list, fetch(socks4Url)...)

	httpURL := NewProxyWeb(name, "https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/http.txt")
	list = append(list, fetch(httpURL)...)
	return list
}
