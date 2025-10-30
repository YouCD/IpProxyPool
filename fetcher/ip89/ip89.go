package ip89

import (
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

func IP89(ctx context.Context) []*database.IP {
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
		}
	}()

	list := make([]*database.IP, 0, 512)
	indexURL := "https://www.89ip.cn"

	doc, _, err := util.Fetch(ctx, indexURL)
	if err != nil {
		log.Errorf("89ip fetch index error: %v", err)
		return list
	}

	pageStr := util.FastText(doc.Find("#layui-laypage-1 > a:nth-child(7)"))
	pageNum, _ := strconv.Atoi(pageStr)
	if pageNum > 5 { // 减少页面数量，防止被反爬和超时
		pageNum = 5
	}

	for i := 1; i <= pageNum; i++ {
		url := fmt.Sprintf("%s/index_%d.html", indexURL, i)
		docPage, _, err := util.Fetch(ctx, url)
		if err != nil {
			log.Errorf("89ip fetch %s error: %v", url, err)
			continue
		}

		docPage.Find("table tbody tr").Each(func(_ int, row *goquery.Selection) {
			ip := util.FastText(row.Find("td:nth-child(1)"))
			port := util.FastText(row.Find("td:nth-child(2)"))
			loc := util.FastText(row.Find("td:nth-child(3)"))

			if ip == "" || port == "" {
				return // 空行跳过
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
		
		// 添加页面间延迟，避免请求过于频繁
		select {
		case <-ctx.Done():
			return list
		case <-time.After(1 * time.Second):
		}
	}
	return list
}