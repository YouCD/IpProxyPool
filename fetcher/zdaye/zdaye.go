package zdaye

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
	"strconv"
	"time"
)

// Zdaye
//
//	@Description: 这个站大爷 只搜索了 https 的代理
func Zdaye() []*database.IP {
	log.Info("[Zdaye] fetch start")
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
		}
	}()
	list := make([]*database.IP, 0)

	indexURL := "https://www.zdaye.com/free/?ip=&adr=&checktime=2&sleep=&cunhuo=&dengji=&nadr=&https=1&yys=&post=&px="
	document, err := fetcher.Fetch(indexURL)
	if err != nil {
		log.Errorf("%s fetch error:%s", indexURL, err)
		return list
	}
	document.Find("table > tbody").Each(func(_ int, selection *goquery.Selection) {
		selection.Find("tr").Each(func(_ int, selection *goquery.Selection) {
			proxyIP := selection.Find("td:nth-child(1)").Text()
			proxyPort := selection.Find("td:nth-child(2)").Text()
			proxyLocation := selection.Find("td:nth-child(4)").Text()
			proxySpeed := selection.Find("td:nth-child(6)").Text()
			ip := new(database.IP)
			ip.ProxyHost = proxyIP
			ip.ProxyPort, _ = strconv.Atoi(proxyPort)
			ip.ProxyType = "https"
			ip.ProxyLocation = proxyLocation
			ip.ProxySpeed, _ = strconv.Atoi(proxySpeed)
			ip.ProxySource = "https://www.zdaye.com"
			ip.CreateTime = time.Now()
			ip.UpdateTime = time.Now()
			list = append(list, ip)
		})
	})
	log.Info("[Zdaye] fetch done")
	return list
}
