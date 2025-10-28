package ip3366

import (
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

//nolint:revive
func Ip3366(ctx context.Context) []*database.IP {
	list := make([]*database.IP, 0)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		// 国内高匿代理
		list = append(list, ip3366(ctx, 1)...)
	}()
	go func() {
		defer wg.Done()
		// 国内普通代理
		list = append(list, ip3366(ctx, 2)...)
	}()

	wg.Wait()

	return list
}

func ip3366(ctx context.Context, proxyType int) []*database.IP {
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
		}
	}()

	list := make([]*database.IP, 0, 1024) // 预分配
	indexURL := "http://www.ip3366.net/free"

	doc, _, err := util.Fetch(ctx, indexURL)
	if err != nil {
		log.Errorf("ip3366 fetch index error: %v", err)
		return list
	}

	pageStr := util.FastText(doc.Find("#listnav > ul > a:nth-child(8)"))
	pageNum, _ := strconv.Atoi(pageStr)
	if pageNum > 50 { // 站点最多 50 页就够了，防止被反爬
		pageNum = 50
	}

	for i := 1; i <= pageNum; i++ {
		url := fmt.Sprintf("%s/?stype=%d&page=%d", indexURL, proxyType, i)
		docPage, _, err := util.Fetch(ctx, url)
		if err != nil {
			log.Errorf("ip3366 fetch %s error: %v", url, err)
			continue
		}

		docPage.Find("table tbody tr").Each(func(_ int, row *goquery.Selection) {
			ip := util.FastText(row.Find("td:nth-child(1)"))
			port := util.FastText(row.Find("td:nth-child(2)"))
			typ := util.FastText(row.Find("td:nth-child(4)"))
			loc := util.FastText(row.Find("td:nth-child(5)"))
			spd := util.FastText(row.Find("td:nth-child(6)"))

			if ip == "" || port == "" {
				return // 空行跳过
			}

			p, _ := strconv.Atoi(port)
			speed, _ := strconv.Atoi(spd)
			list = append(list, &database.IP{
				ProxyHost:     ip,
				ProxyPort:     p,
				ProxyType:     typ,
				ProxyLocation: loc,
				ProxySpeed:    speed,
				ProxySource:   indexURL,
				CreateTime:    time.Now(),
				UpdateTime:    time.Now(),
			})
		})
	}
	return list
}
