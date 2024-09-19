package github

import (
	"IpProxyPool/middleware/database"
)

func Anonym0usWork1221() []*database.IP {
	list := make([]*database.IP, 0)
	name := "Anonym0usWork1221"

	// http
	http := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/http_proxies.txt"))
	list = append(list, http...)

	// https
	https := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/https_proxies.txt"))
	list = append(list, https...)
	// socks4
	socks4 := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/socks4_proxies.txt"))
	list = append(list, socks4...)
	// socks5
	socks5 := fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/socks5_proxies.txt"))
	list = append(list, socks5...)
	return list
}
