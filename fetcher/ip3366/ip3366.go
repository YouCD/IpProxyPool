package ip3366

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
	"strconv"
	"strings"
)

func Ip3366() []*database.IP {
	list := make([]*database.IP, 0)
	// 国内高匿代理
	list = append(list, ip3366(1)...)
	// 国内普通代理
	list = append(list, ip3366(2)...)
	return list
}

func ip3366(proxyType int) []*database.IP {
	log.Info("[ip3366] fetch start")
	defer func() {
		recover()
		log.Warn("[ip3366] fetch error")
	}()
	list := make([]*database.IP, 0)

	indexUrl := "http://www.ip3366.net/free"
	document, err := fetcher.Fetch(indexUrl)
	if err != nil {
		log.Errorf("ip3366 fetch index error:%s", err)
		return list
	}
	pageNum := document.Find("#listnav > ul > a:nth-child(8)").Text()
	num, _ := strconv.Atoi(pageNum)
	for i := 1; i <= num; i++ {
		url := fmt.Sprintf("%s/?stype=%d&page=%d", indexUrl, proxyType, i)
		documentA, err := fetcher.Fetch(url)
		if err != nil {
			log.Errorf("%s documentA error:%s", url, err)
			continue
		}
		documentA.Find("table > tbody").Each(func(i int, selection *goquery.Selection) {
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
				ip.ProxySource = "http://www.ip3366.net"
				ip.CreateTime = util.FormatDateTime()
				ip.UpdateTime = util.FormatDateTime()
				list = append(list, ip)
			})
		})
	}
	log.Info("[ip3366] fetch done")
	return list
}
