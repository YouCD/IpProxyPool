package proxylistplus

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
	"time"

	"strconv"
	"strings"
)

func ProxyListPlus() []*database.IP {
	list := make([]*database.IP, 0)
	indexURL := "https://list.proxylistplus.com"
	for i := 1; i <= 6; i++ {
		url := fmt.Sprintf("%s/Fresh-HTTP-Proxy-List-%d", indexURL, i)
		document, err := fetcher.Fetch(url)
		if err != nil {
			log.Errorf("[proxylistplus] document failed,err:%s", err)
			return list
		}
		document.Find("table.bg > tbody").Each(func(_ int, selection *goquery.Selection) {
			selection.Find("tr").Each(func(_ int, selection *goquery.Selection) {
				proxyIP := strings.TrimSpace(selection.Find("td:nth-child(2)").Text())
				proxyPort := strings.TrimSpace(selection.Find("td:nth-child(3)").Text())
				proxyLocation := strings.TrimSpace(selection.Find("td:nth-child(5)").Text())

				ip := new(database.IP)
				ip.ProxyHost = proxyIP
				ip.ProxyPort, _ = strconv.Atoi(proxyPort)
				ip.ProxyType = "http"
				ip.ProxyLocation = proxyLocation
				ip.ProxySpeed = 100
				ip.ProxySource = "https://list.proxylistplus.com"
				ip.CreateTime = time.Now()
				ip.UpdateTime = time.Now()
				list = append(list, ip)
			})
		})
	}
	return list
}
