package github

import (
	"IpProxyPool/middleware/database"
)

func ProxyScraper() []*database.IP {
	list := make([]*database.IP, 0)
	name := "ProxyScraper"

	// http
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/http.txt"))...)
	// socks4
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/socks4.txt"))...)
	// socks5
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/socks5.txt"))...)
	return list
}
