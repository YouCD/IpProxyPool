package github

import (
	"IpProxyPool/middleware/database"
)

func Anonym0usWork1221() []*database.IP {
	list := make([]*database.IP, 0)
	name := "Anonym0usWork1221"

	// http
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/http_proxies.txt"))...)

	// https
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/https_proxies.txt"))...)
	// socks4
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/socks4_proxies.txt"))...)
	// socks5
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/socks5_proxies.txt"))...)
	return list
}
