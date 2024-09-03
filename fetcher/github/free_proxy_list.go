package github

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"github.com/youcd/toolkit/log"
	"strconv"
	"strings"
)

func FreeProxyList() []*database.IP {
	list := make([]*database.IP, 0)
	socks5Url := setProxyWeb("https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks5/data.txt")
	list = append(list, freeProxyListFetch(socks5Url, "socks5://")...)

	socks4Url := setProxyWeb("https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks4/data.txt")
	list = append(list, freeProxyListFetch(socks4Url, "socks4://")...)

	httpUrl := setProxyWeb("https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/http/data.txt")
	list = append(list, freeProxyListFetch(httpUrl, "http://")...)

	return list
}

func freeProxyListFetch(urlStr, replaceStr string) []*database.IP {
	list := make([]*database.IP, 0)
	log.Infof("[FreeProxyList] fetch start: %s", urlStr)
	document, err := fetcher.Fetch(urlStr)
	if err != nil {
		log.Errorf("%s fetch failed,err: %s", urlStr, err)
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
		ip.CreateTime = util.FormatDateTime()
		ip.UpdateTime = util.FormatDateTime()
		list = append(list, ip)
	}
	log.Infof("[FreeProxyList] fetch done: %s, count: %d", urlStr, len(list))
	return list
}
