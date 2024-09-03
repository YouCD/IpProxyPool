package github

import "IpProxyPool/middleware/database"

func FreshProxyList() []*database.IP {
	list := make([]*database.IP, 0)
	httpsUrl := setProxyWeb("https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/https.txt")
	list = append(list, fetch("FreshProxyList", httpsUrl)...)

	socks5Url := setProxyWeb("https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/socks5.txt")
	list = append(list, fetch("FreshProxyList", socks5Url)...)

	socks4Url := setProxyWeb("https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/socks4.txt")
	list = append(list, fetch("FreshProxyList", socks4Url)...)

	httpUrl := setProxyWeb("https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/http.txt")
	list = append(list, fetch("FreshProxyList", httpUrl)...)

	return list
}
