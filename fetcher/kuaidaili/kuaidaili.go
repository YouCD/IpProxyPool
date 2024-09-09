package kuaidaili

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
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
	indexURL := "https://www.kuaidaili.com/free"
	document, err := fetcher.Fetch(indexURL)
	if err != nil {
		log.Warnf("KuaiDaiLi: 类型: %s  err:%s", proxyTypeStr(proxyType), err)
		return list
	}
	pageNum := document.Find("#listnav > ul > li:nth-child(9) > a").Text()
	num, _ := strconv.Atoi(pageNum)
	for i := 1; i <= num; i++ {
		//  休眠3秒，防止被封
		time.Sleep(3 * time.Second)
		url := fmt.Sprintf("%s/%s/%d", indexURL, proxyType, i)

		documentA, err := fetcher.Fetch(url)
		if err != nil {
			log.Errorf("KuaiDaiLi: 类型:%s err:%s", proxyTypeStr(proxyType), err)
			continue
		}

		documentA.Find("table > tbody").Each(func(_ int, selection *goquery.Selection) {
			selection.Find("tr").Each(func(_ int, selection *goquery.Selection) {
				proxyIP := strings.TrimSpace(selection.Find("td:nth-child(1)").Text())
				proxyPort := strings.TrimSpace(selection.Find("td:nth-child(2)").Text())
				proxyTyp := strings.TrimSpace(selection.Find("td:nth-child(4)").Text())
				proxyLocation := strings.TrimSpace(selection.Find("td:nth-child(5)").Text())
				proxySpeed := strings.TrimSpace(selection.Find("td:nth-child(6)").Text())

				ip := new(database.IP)
				ip.ProxyHost = proxyIP
				ip.ProxyPort, _ = strconv.Atoi(proxyPort)
				ip.ProxyType = proxyTyp
				ip.ProxyLocation = proxyLocation
				ip.ProxySpeed, _ = strconv.Atoi(proxySpeed)
				ip.ProxySource = "https://www.kuaidaili.com"
				ip.CreateTime = time.Now()
				ip.UpdateTime = time.Now()
				list = append(list, ip)
			})
		})
	}
	log.Info("KuaiDaiLi fetch done")
	return list
}
