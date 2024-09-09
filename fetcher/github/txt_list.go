package github

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"fmt"
	"github.com/youcd/toolkit/log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

func fetch(name, urlStr string) []*database.IP {
	log.Infof("[%s] fetch start", name)
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("[%s] fetch error:%s", name, r)
		}
	}()
	list := make([]*database.IP, 0)
	parse, err := url.Parse(urlStr)
	if err != nil {
		log.Error(err)
		return nil
	}

	proxySource := fmt.Sprintf("%s://%s", parse.Scheme, parse.Host)
	document, err := fetcher.Fetch(urlStr)
	if err != nil {
		log.Errorf("[%s] fetch failed,url: %s,err:%s", name, urlStr, err)
		return list
	}
	split := strings.Split(document.Text(), "\n")
	var wg sync.WaitGroup
	for _, address := range split {
		if address == "" {
			continue
		}
		wg.Add(1)
		go func(ipPort string) {
			defer wg.Done()
			if _, err := net.DialTimeout("tcp", ipPort, 3*time.Second); err != nil {
				return
			}
			ipPortObj := strings.Split(ipPort, ":")
			ip := new(database.IP)
			ip.ProxyHost = ipPortObj[0]
			ip.ProxyPort, _ = strconv.Atoi(ipPortObj[1])
			ip.ProxyLocation = parse.Host
			ip.ProxySpeed = 100
			ip.ProxySource = proxySource
			ip.CreateTime = time.Now()
			ip.UpdateTime = time.Now()
			list = append(list, ip)
		}(address)
	}
	wg.Wait()

	log.Infof("[%s] fetch done", name)
	return list
}
