package ip66

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

//nolint:revive
func Ip66() []*database.IP {
	log.Info("[66ip] fetch start")
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
		}
	}()
	list := make([]*database.IP, 0)

	indexUrl := "http://www.66ip.cn"
	for i := 1; i <= 100; i++ {
		url := fmt.Sprintf("%s/%d.html", indexUrl, i)
		document, err := fetcher.Fetch(url)
		if err != nil {
			log.Errorf("%s fetch failed,err: %s", indexUrl, err)
			return list
		}
		document.Find("table > tbody").Each(func(i int, selection *goquery.Selection) {
			selection.Find("tr").NextAll().Each(func(i int, selection *goquery.Selection) {
				proxyIp := strings.TrimSpace(selection.Find("td:nth-child(1)").Text())
				proxyPort := strings.TrimSpace(selection.Find("td:nth-child(2)").Text())
				proxyLocation := strings.TrimSpace(selection.Find("td:nth-child(3)").Text())

				ip := new(database.IP)
				ip.ProxyHost = proxyIp
				ip.ProxyPort, _ = strconv.Atoi(proxyPort)
				ip.ProxyType = "http"
				ip.ProxyLocation = proxyLocation
				ip.ProxySpeed = 100
				ip.ProxySource = "http://www.66ip.cn"
				ip.CreateTime = time.Now()
				ip.UpdateTime = time.Now()
				list = append(list, ip)
			})
		})
	}
	log.Info("[66ip] fetch done")
	return list
}
