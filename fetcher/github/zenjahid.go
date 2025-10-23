package github

import (
	"IpProxyPool/middleware/database"
)

func Zenjahid() []*database.IP {
	list := make([]*database.IP, 0)
	name := "zenjahid"
	// http
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/zenjahid/FreeProxy4u/main/http.txt"))...)
	// socks4
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/zenjahid/FreeProxy4u/main/socks4.txt"))...)
	// socks5
	list = append(list, fetch(NewProxyWeb(name, "https://raw.githubusercontent.com/zenjahid/FreeProxy4u/main/socks5.txt"))...)

	return list
}
