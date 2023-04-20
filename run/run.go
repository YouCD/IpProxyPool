package run

import (
	"IpProxyPool/fetcher/ip3366"
	"IpProxyPool/fetcher/ip66"
	"IpProxyPool/fetcher/ip89"
	"IpProxyPool/middleware/database"
	"IpProxyPool/middleware/storage"
	logger "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func Task() {
	ipChan := make(chan *database.IP, 2000)

	// Check the IPs in DB
	go func() {
		storage.CheckProxyDB()
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
		logger.Printf("Chan: %v, IP: %d\n", len(ipChan), nums)
		if len(ipChan) < 100 {
			go run(ipChan)
		}
		time.Sleep(300 * time.Second)
	}
}

func run(ipChan chan<- *database.IP) {
	var wg sync.WaitGroup
	siteFuncList := []func() []*database.IP{
		ip66.Ip66,
		ip89.Ip89,
		ip3366.Ip33661,
		ip3366.Ip33662,
		//kuaidaili.KuaiDaiLiInha,
		//kuaidaili.KuaiDaiLiIntr,
		//proxylistplus.ProxyListPlus,
	}
	for _, siteFunc := range siteFuncList {
		wg.Add(1)
		go func(siteFunc func() []*database.IP) {
			temp := siteFunc()
			for _, v := range temp {
				ipChan <- v
			}
			wg.Done()
		}(siteFunc)
	}
	wg.Wait()
	logger.Println("All getters finished.")
}
