package run

import (
	"IpProxyPool/fetcher/geonode"
	"IpProxyPool/fetcher/github"
	"IpProxyPool/fetcher/ip3366"
	"IpProxyPool/fetcher/ip89"
	"IpProxyPool/fetcher/proxylistplus"
	"IpProxyPool/middleware/database"
	"IpProxyPool/middleware/storage"
	"context"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
)

func Task(ctx context.Context) {
	ipChan := make(chan *database.IP, 2000)

	// 循环检测数据库中的IP
	go func() {
		c := cron.New()
		_, _ = c.AddFunc("*/30 * * * *", func() {
			storage.CheckProxyDB(ctx)
		})
		c.Start()
	}()

	// Check the IPs in channel
	numConsumers := 30 // 设置消费者数量
	log.Debugf("Starting consumer total %d", numConsumers)
	for i := range numConsumers {
		go func(consumerID int) {
			for {
				select {
				case <-ctx.Done():
					log.Infof("Consumer %d stopping due to context cancellation", consumerID)
					return
				case ip := <-ipChan:
					if ip == nil {
						log.Warnf("Consumer %d received nil IP, skipping...", consumerID)
						continue
					}
					log.Infow("CheckProxy", "consumerID", consumerID, "ipChan len", len(ipChan), "msg", storage.CheckProxy(ctx, ip))
				}
			}
		}(i)
	}

	go func() {
		c := cron.New()
		_, _ = c.AddFunc("*/1 * * * *", func() {
			nums := database.CountIP()
			log.Infof("count for Chan: %v, count for database : %d", len(ipChan), nums)
			run(ctx, ipChan)
		})
		c.Start()
	}()
}

func run(ctx context.Context, ipChan chan<- *database.IP) {
	var wg sync.WaitGroup

	type fetcher func(ctx context.Context) []*database.IP
	siteFuncList := map[string]fetcher{
		"89ip":              ip89.IP89,
		"ip3366":            ip3366.Ip3366,
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
		"R00tee":            github.R00tee,
	}
	// --- 2. 设置上下文 + 控制参数 ---
	for name, siteFunc := range siteFuncList {
		wg.Add(1)
		go func(name string, fetcherFunc fetcher) {
			timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer func() {
				cancel()
				wg.Done()
			}()
			temp := fetcherFunc(timeoutCtx)
			log.Infof("[%s] Get IP: %d", name, len(temp))
			for _, ip := range temp {
				select {
				case <-ctx.Done():
					log.Warnf("[%s] canceled before sending ip", name)
					return
				case ipChan <- ip:
				case <-time.After(2 * time.Second):
					log.Warnf("send timeout, channel likely full: %s, len: %d", "ipChan", len(ipChan))
				}
			}
		}(name, siteFunc)
	}
	wg.Wait()
	log.Info("All getters finished.")
}
