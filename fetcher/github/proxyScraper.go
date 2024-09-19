package github

import (
	"IpProxyPool/middleware/database"
)

func ProxyScraper() []*database.IP {
	list := make([]*database.IP, 0)
	name := "ProxyScraper"

	// http
	http := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/http.txt"))
	list = append(list, http...)
	// socks4
	socks4 := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/socks4.txt"))
	list = append(list, socks4...)
	// socks5
	socks5 := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/socks5.txt"))
	list = append(list, socks5...)
	return list
}
