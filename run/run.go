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
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
	"sync"
)

func Task() {
	ipChan := make(chan *database.IP, 2000)

	// 循环检测数据库中的IP
	go func() {
		c := cron.New()
		_, _ = c.AddFunc("*/30 * * * *", func() {
			storage.CheckProxyDB()
		})
		c.Start()
	}()

	// Check the IPs in channel
	numConsumers := 30 // 设置消费者数量
	for i := range numConsumers {
		go func(consumerID int) {
			log.Infof("Starting consumer %d", consumerID)
			for {
				ip := <-ipChan
				if ip == nil {
					log.Warnf("Consumer %d received nil IP, skipping...", consumerID)
					continue
				}
				proxyStr := fmt.Sprintf("%s:%d", ip.ProxyHost, ip.ProxyPort)
				log.Infof("Consumer %d checking IP: %s", consumerID, proxyStr)
				storage.CheckProxy(ip)
			}
		}(i)
	}

	go func() {
		c := cron.New()
		_, _ = c.AddFunc("*/5 * * * *", func() {
			nums := database.CountIP()
			log.Infof("count for Chan: %v, count for database : %d", len(ipChan), nums)
			run(ipChan)
		})
		c.Start()
	}()
}

func run(ipChan chan<- *database.IP) {
	var wg sync.WaitGroup

	type fetcher func() []*database.IP
	siteFuncList := map[string]fetcher{
		// "66ip":          ip66.Ip66,
		"89ip":              ip89.Ip89,
		"ip3366":            ip3366.Ip3366,
		"Zdaye":             zdaye.Zdaye,
		"KuaiDaiLi":         kuaidaili.KuaiDaiLi,
		"proxylistplus":     proxylistplus.ProxyListPlus,
		"TheSpeedX":         github.TheSpeedX,
		"OpenProxyList":     github.OpenProxyList,
		"Geonode":           geonode.Geonode,
		"ZloiUser":          github.ZloiUser,
		"FreeProxyList":     github.FreeProxyList,
		"Vakhov":            github.Vakhov,
		"Yemixzy":           github.Yemixzy,
		"Zaeem20":           github.Zaeem20,
		"Anonym0usWork1221": github.Anonym0usWork1221,
		"Zenjahid":          github.Zenjahid,
		"ProxyScraper":      github.ProxyScraper,
	}

	for name, siteFunc := range siteFuncList {
		wg.Add(1)
		go func(name string, fetcherFunc fetcher) {
			defer wg.Done()
			temp := fetcherFunc()
			log.Infof("[%s] Get IP: %d", name, len(temp))
			for _, ip := range temp {
				proxyStr := fmt.Sprintf("%s:%d", ip.ProxyHost, ip.ProxyPort)
				log.Debugf("[%s] Send proxy: %s", name, proxyStr)
				ipChan <- ip
				log.Debugf("[%s] Send OK: %s", name, proxyStr)
			}
		}(name, siteFunc)
	}
	wg.Wait()
	log.Info("All getters finished.")
}
