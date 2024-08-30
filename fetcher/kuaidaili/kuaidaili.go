package kuaidaili

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
	"strconv"
	"strings"
	"time"
)

func KuaiDaiLi() []*database.IP {
	list := make([]*database.IP, 0)
	// 国内高匿代理
	list = append(list, kuaiDaiLi("inha")...)
	// 国内普通代理
	list = append(list, kuaiDaiLi("intr")...)
	return list
}

func proxyTypeStr(typ string) string {
	switch typ {
	case "inha":
		return "国内高匿代理"
	case "intr":
		return "国内普通代理"
	}
	return "未知"
}
func kuaiDaiLi(proxyType string) []*database.IP {
	log.Infow("KuaiDaiLi", "类型", proxyTypeStr(proxyType))

	list := make([]*database.IP, 0)

	indexUrl := "https://www.kuaidaili.com/free"
	fetchIndex := fetcher.Fetch(indexUrl)
	if fetchIndex == nil {
		log.Warnf("KuaiDaiLi: 类型:%s  限流", proxyTypeStr(proxyType))
		return nil
	}
	pageNum := fetchIndex.Find("#listnav > ul > li:nth-child(9) > a").Text()
	num, _ := strconv.Atoi(pageNum)
	for i := 1; i <= num; i++ {
		//  休眠3秒，防止被封
		time.Sleep(3 * time.Second)
		url := fmt.Sprintf("%s/%s/%d", indexUrl, proxyType, i)

		fetch := fetcher.Fetch(url)
		if fetch == nil {
			log.Warnw("KuaiDaiLi", "类型", proxyTypeStr(proxyType))
			continue
		}

		fetch.Find("table > tbody").Each(func(i int, selection *goquery.Selection) {
			selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
				proxyIp := strings.TrimSpace(selection.Find("td:nth-child(1)").Text())
				proxyPort := strings.TrimSpace(selection.Find("td:nth-child(2)").Text())
				proxyType := strings.TrimSpace(selection.Find("td:nth-child(4)").Text())
				proxyLocation := strings.TrimSpace(selection.Find("td:nth-child(5)").Text())
				proxySpeed := strings.TrimSpace(selection.Find("td:nth-child(6)").Text())

				ip := new(database.IP)
				ip.ProxyHost = proxyIp
				ip.ProxyPort, _ = strconv.Atoi(proxyPort)
				ip.ProxyType = proxyType
				ip.ProxyLocation = proxyLocation
				ip.ProxySpeed, _ = strconv.Atoi(proxySpeed)
				ip.ProxySource = "https://www.kuaidaili.com"
				ip.CreateTime = util.FormatDateTime()
				ip.UpdateTime = util.FormatDateTime()
				list = append(list, ip)
			})
		})
	}
	log.Info("KuaiDaiLi fetch done")
	return list
}
