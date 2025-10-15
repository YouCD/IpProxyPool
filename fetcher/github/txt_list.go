package github

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/youcd/toolkit/log"
)

func fetch(proxyWeb *ProxyWeb) []*database.IP {
	if proxyWeb == nil {
		return nil
	}

	var count int
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("[%s] fetch error:%s", proxyWeb.Name, r)
		}
	}()
	list := make([]*database.IP, 0)
Retry:
	document, err := fetcher.Fetch(proxyWeb.GetFullURL())
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		proxyWeb.ChangeProxy()
		count++
		if count < 3 {
			log.Errorf("[%s] ChangeProxy: %s ", proxyWeb.Name, proxyWeb.ProxyURL)
			goto Retry
		}
		log.Errorf("[%s] fetch failed,url: %s,err:%s", proxyWeb.Name, proxyWeb.GetFullURL(), err)
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
			ip.ProxyLocation = proxyWeb.Name
			ip.ProxySpeed = 100
			ip.ProxySource = proxyWeb.Name
			ip.CreateTime = time.Now()
			ip.UpdateTime = time.Now()
			list = append(list, ip)
		}(address)
	}
	wg.Wait()
	return list
}
