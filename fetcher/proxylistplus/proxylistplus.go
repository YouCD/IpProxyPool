package proxylistplus

import (
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"fmt"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

func ProxyListPlus() []*database.IP {
	list := make([]*database.IP, 0, 512) // 预分配
	indexURL := "https://list.proxylistplus.com"
	for i := 1; i <= 6; i++ {
		url := fmt.Sprintf("%s/Fresh-HTTP-Proxy-List-%d", indexURL, i)
		doc, _, err := util.Fetch(url)
		if err != nil {
			log.Errorf("[proxylistplus] fetch failed: %v", err)
			continue // 不要 return，跳过即可
		}

		doc.Find("table.bg tbody tr").Each(func(_ int, row *goquery.Selection) {
			// 只取 3 个 td，不整行 Text()
			ip := util.FastText(row.Find("td:nth-child(2)"))
			port := util.FastText(row.Find("td:nth-child(3)"))
			loc := util.FastText(row.Find("td:nth-child(5)"))

			if ip == "" || port == "" {
				return // 跳过空行
			}
			p, _ := strconv.Atoi(port)
			list = append(list, &database.IP{
				ProxyHost:     ip,
				ProxyPort:     p,
				ProxyType:     "http",
				ProxyLocation: loc,
				ProxySpeed:    100,
				ProxySource:   indexURL,
				CreateTime:    time.Now(),
				UpdateTime:    time.Now(),
			})
		})
	}
	return list
}
