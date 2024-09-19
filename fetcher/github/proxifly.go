package github

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"github.com/youcd/toolkit/log"
	"strconv"
	"strings"
	"time"
)

func FreeProxyList() []*database.IP {
	list := make([]*database.IP, 0)
	name := "FreeProxyList"
	socks5Url := NewProxyWeb(name, "https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks5/data.txt")
	list = append(list, freeProxyListFetch(socks5Url, "socks5://")...)

	socks4Url := NewProxyWeb(name, "https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks4/data.txt")
	list = append(list, freeProxyListFetch(socks4Url, "socks4://")...)

	httpURL := NewProxyWeb(name, "https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/http/data.txt")
	list = append(list, freeProxyListFetch(httpURL, "http://")...)
	return list
}

func freeProxyListFetch(urlStr *ProxyWeb, replaceStr string) []*database.IP {
	list := make([]*database.IP, 0)
	document, err := fetcher.Fetch(urlStr.GetFullURL())
	if err != nil {
		log.Errorf("%s fetch failed,err: %s", urlStr.Name, err)
		return list
	}

	for _, s := range strings.Split(document.Text(), "\n") {
		s := strings.ReplaceAll(s, replaceStr, "")
		split := strings.Split(s, ":")
		if len(split) < 2 {
			continue
		}
		ip := new(database.IP)
		ip.ProxyHost = split[0]
		ip.ProxyPort, _ = strconv.Atoi(split[1])
		ip.ProxyLocation = "free-proxy-list"
		ip.ProxySpeed = 100
		ip.ProxySource = "https://github.com/proxifly/free-proxy-list"
		ip.CreateTime = time.Now()
		ip.UpdateTime = time.Now()
		list = append(list, ip)
	}
	return list
}
