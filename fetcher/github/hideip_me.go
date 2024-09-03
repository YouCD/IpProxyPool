package github

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"github.com/youcd/toolkit/log"
	"strconv"
	"strings"
)

func HideIPMe() []*database.IP {
	list := make([]*database.IP, 0)

	HTTPSUrl := setProxyWeb("https://raw.githubusercontent.com/zloi-user/hideip.me/main/https.txt")
	list = append(list, hideIPMeFetch(HTTPSUrl)...)

	socks4Url := setProxyWeb("https://raw.githubusercontent.com/zloi-user/hideip.me/main/socks4.txt")
	list = append(list, hideIPMeFetch(socks4Url)...)

	socks5 := setProxyWeb("https://raw.githubusercontent.com/zloi-user/hideip.me/main/socks5.txt")
	list = append(list, hideIPMeFetch(socks5)...)

	return list
}
func hideIPMeFetch(urlStr string) []*database.IP {
	log.Infof("[hideip.me] fetch start: %s", urlStr)
	list := make([]*database.IP, 0)
	document, err := fetcher.Fetch(urlStr)
	if err != nil {
		log.Errorf("%s fetch failed,err:%s", urlStr, err)
		return list
	}
	for _, s := range strings.Split(document.Text(), "\n") {
		split := strings.Split(s, ":")
		if len(split) < 3 {
			continue
		}
		ip := new(database.IP)
		ip.ProxyHost = split[0]
		ip.ProxyPort, _ = strconv.Atoi(split[1])
		ip.ProxyLocation = split[2]
		ip.ProxySpeed = 100
		ip.ProxySource = "https://github.com/zloi-user/hideip.me"
		ip.CreateTime = util.FormatDateTime()
		ip.UpdateTime = util.FormatDateTime()
		list = append(list, ip)
	}
	log.Infof("[hideip.me] fetch done: %s, count: %d", urlStr, len(list))
	return list
}
