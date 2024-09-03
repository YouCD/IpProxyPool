package ip89

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

func Ip89() []*database.IP {
	log.Info("[89ip] fetch start")
	defer func() {
		recover()
		log.Warn("[89ip] fetch error")
	}()
	list := make([]*database.IP, 0)

	indexUrl := "https://www.89ip.cn/"
	document, err := fetcher.Fetch(indexUrl)
	if err != nil {
		log.Errorf("%s fetch error:%s", indexUrl, err)
		return list
	}
	pageNum := document.Find("#layui-laypage-1 > a:nth-child(7)").Text()
	num, _ := strconv.Atoi(pageNum)
	for i := 1; i <= num; i++ {
		url := fmt.Sprintf("%s/index_%d.html", indexUrl, i)
		documentA, err := fetcher.Fetch(url)
		if err != nil {
			log.Errorf("%s document error:%s", indexUrl, err)
			continue
		}
		documentA.Find("table > tbody").Each(func(i int, selection *goquery.Selection) {
			selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
				proxyIp := strings.TrimSpace(selection.Find("td:nth-child(1)").Text())
				proxyPort := strings.TrimSpace(selection.Find("td:nth-child(2)").Text())
				proxyLocation := strings.TrimSpace(selection.Find("td:nth-child(3)").Text())

				ip := new(database.IP)
				ip.ProxyHost = proxyIp
				ip.ProxyPort, _ = strconv.Atoi(proxyPort)
				ip.ProxyType = "http"
				ip.ProxyLocation = proxyLocation
				ip.ProxySpeed = 100
				ip.ProxySource = "https://www.89ip.cn"
				ip.CreateTime = util.FormatDateTime()
				ip.UpdateTime = util.FormatDateTime()
				list = append(list, ip)
			})
		})
	}
	log.Info("[89ip] fetch done")
	return list
}
