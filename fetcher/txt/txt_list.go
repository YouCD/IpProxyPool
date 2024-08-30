package txt

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"fmt"
	"github.com/youcd/toolkit/log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ProxyType string

func OpenProxyList() []*database.IP {
	list := make([]*database.IP, 0)
	// http
	http := fetch("OpenProxyList", "https://api.openproxylist.xyz/http.txt")
	list = append(list, http...)
	// socks5
	socks5 := fetch("OpenProxyList", "https://api.openproxylist.xyz/socks5.txt")
	list = append(list, socks5...)
	//	 socks4
	socks4 := fetch("OpenProxyList", "https://api.openproxylist.xyz/socks4.txt")
	list = append(list, socks4...)

	return list
}

func TheSpeedX() []*database.IP {
	list := make([]*database.IP, 0)
	// http
	httpIps := fetch("ProxyList", "https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt")
	list = append(list, httpIps...)
	socks4Ips := fetch("ProxyList", "https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks4.txt")
	list = append(list, socks4Ips...)
	socks5Ips := fetch("ProxyList", "https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks5.txt")
	list = append(list, socks5Ips...)
	return list
}

func fetch(name, urlStr string) []*database.IP {
	log.Infof("[%s] fetch start", name)
	defer func() {
		recover()
		log.Warnf("[%s] fetch error", name)
	}()
	list := make([]*database.IP, 0)
	parse, err := url.Parse(urlStr)
	if err != nil {
		log.Error(err)
		return nil
	}

	proxySource := fmt.Sprintf("%s://%s", parse.Scheme, parse.Host)

	fetchIndex := fetcher.Fetch(urlStr)
	split := strings.Split(fetchIndex.Text(), "\n")
	var wg sync.WaitGroup
	for _, address := range split {
		if address == "" {
			continue
		}
		wg.Add(1)
		go func(ipPort string) {
			defer wg.Done()
			_, err := net.DialTimeout("tcp", ipPort, 3*time.Second)
			if err != nil {
				return
			} else {
				ipPortObj := strings.Split(ipPort, ":")
				ip := new(database.IP)
				ip.ProxyHost = ipPortObj[0]
				ip.ProxyPort, _ = strconv.Atoi(ipPortObj[1])
				ip.ProxyLocation = parse.Host
				ip.ProxySpeed = 100
				ip.ProxySource = proxySource
				ip.CreateTime = util.FormatDateTime()
				ip.UpdateTime = util.FormatDateTime()
				list = append(list, ip)
			}
		}(address)
	}
	wg.Wait()

	log.Infof("[%s] fetch done", name)
	return list
}
