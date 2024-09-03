package proxylistplus

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

func ProxyListPlus() []*database.IP {
	log.Info("[proxylistplus] fetch start")
	list := make([]*database.IP, 0)
	indexUrl := "https://list.proxylistplus.com"
	for i := 1; i <= 6; i++ {
		url := fmt.Sprintf("%s/Fresh-HTTP-Proxy-List-%d", indexUrl, i)
		document, err := fetcher.Fetch(url)
		if err != nil {
			log.Errorf("[proxylistplus] document failed,err:%s", err)
			return list
		}
		document.Find("table.bg > tbody").Each(func(i int, selection *goquery.Selection) {
			selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
				proxyIp := strings.TrimSpace(selection.Find("td:nth-child(2)").Text())
				proxyPort := strings.TrimSpace(selection.Find("td:nth-child(3)").Text())
				proxyLocation := strings.TrimSpace(selection.Find("td:nth-child(5)").Text())

				ip := new(database.IP)
				ip.ProxyHost = proxyIp
				ip.ProxyPort, _ = strconv.Atoi(proxyPort)
				ip.ProxyType = "http"
				ip.ProxyLocation = proxyLocation
				ip.ProxySpeed = 100
				ip.ProxySource = "https://list.proxylistplus.com"
				ip.CreateTime = util.FormatDateTime()
				ip.UpdateTime = util.FormatDateTime()
				list = append(list, ip)
			})
		})
	}
	log.Info("[proxylistplus] fetch done")
	return list
}
