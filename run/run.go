package run

import (
	"IpProxyPool/fetcher/geonode"
	"IpProxyPool/fetcher/github"
	"IpProxyPool/fetcher/ip3366"
	"IpProxyPool/fetcher/ip89"
	"IpProxyPool/fetcher/kuaidaili"
	"IpProxyPool/fetcher/proxylistplus"
	"IpProxyPool/fetcher/zdaye"
	"IpProxyPool/middleware/database"
	"IpProxyPool/middleware/storage"
	"github.com/youcd/toolkit/log"
	"sync"
	"time"
)

func Task() {
	ipChan := make(chan *database.IP, 2000)

	// 循环检测数据库中的IP
	go func() {
		for {
			log.Info("Checking IPs in DB...")
			storage.CheckProxyDB()
			time.Sleep(10 * time.Minute)
		}
	}()

	// Check the IPs in channel
	for i := 0; i < 50; i++ {
		go func() {
			for {
				storage.CheckProxy(<-ipChan)
			}
		}()
	}

	// Start getters to scraper IP and put it in channel
	for {
		nums := database.CountIp()
		log.Infof("Chan: %v, IP: %d", len(ipChan), nums)
		if len(ipChan) < 100 {
			go run(ipChan)
		}
		time.Sleep(300 * time.Second)
	}
}

func run(ipChan chan<- *database.IP) {
	var wg sync.WaitGroup

	type fetcher func() []*database.IP
	siteFuncList := map[string]fetcher{
		//"66ip":          ip66.Ip66,
		"89ip":           ip89.Ip89,
		"ip3366":         ip3366.Ip3366,
		"站大爷":            zdaye.Zdaye,
		"快代理":            kuaidaili.KuaiDaiLi,
		"proxylistplus":  proxylistplus.ProxyListPlus,
		"TheSpeedX":      github.TheSpeedX,
		"OpenProxyList":  github.OpenProxyList,
		"Geonode":        geonode.Geonode,
		"HideIPMe":       github.HideIPMe,
		"FreeProxyList":  github.FreeProxyList,
		"FreshProxyList": github.FreshProxyList,
	}

	for name, siteFunc := range siteFuncList {
		wg.Add(1)
		go func(siteFunc fetcher) {
			temp := siteFunc()
			log.Infof("Get %d IP from %s", len(temp), name)
			for _, v := range temp {
				ipChan <- v
			}
			wg.Done()
		}(siteFunc)
	}
	wg.Wait()
	log.Info("All getters finished.")
}
