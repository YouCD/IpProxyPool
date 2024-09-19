package github

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"github.com/youcd/toolkit/log"
	"strconv"
	"strings"
	"time"
)

func ZloiUser() []*database.IP {
	list := make([]*database.IP, 0)

	name := "zloiUser"
	HTTPSUrl := NewProxyWeb(name, "https://raw.githubusercontent.com/zloi-user/hideip.me/main/https.txt")
	list = append(list, hideIPMeFetch(HTTPSUrl)...)

	socks4Url := NewProxyWeb(name, "https://raw.githubusercontent.com/zloi-user/hideip.me/main/socks4.txt")
	list = append(list, hideIPMeFetch(socks4Url)...)

	socks5 := NewProxyWeb(name, "https://raw.githubusercontent.com/zloi-user/hideip.me/main/socks5.txt")
	list = append(list, hideIPMeFetch(socks5)...)
	return list
}
func hideIPMeFetch(urlStr *ProxyWeb) []*database.IP {
	list := make([]*database.IP, 0)
	document, err := fetcher.Fetch(urlStr.GetFullURL())
	if err != nil {
		log.Errorf("%s fetch failed,err:%s", urlStr.GetFullURL(), err)
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
		ip.CreateTime = time.Now()
		ip.UpdateTime = time.Now()
		list = append(list, ip)
	}
	return list
}
