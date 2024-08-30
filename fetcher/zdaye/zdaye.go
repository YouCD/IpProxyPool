package zdaye

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
	"strconv"
)

// Zdaye
//
//	@Description: 这个站大爷 只搜索了 https 的代理
func Zdaye() []*database.IP {

	log.Info("[Zdaye] fetch start")
	defer func() {
		recover()
		log.Warn("[Zdaye] fetch error")
	}()
	list := make([]*database.IP, 0)

	indexUrl := "https://www.zdaye.com/free/?ip=&adr=&checktime=2&sleep=&cunhuo=&dengji=&nadr=&https=1&yys=&post=&px="
	fetchIndex := fetcher.Fetch(indexUrl)
	fetchIndex.Find("table > tbody").Each(func(i int, selection *goquery.Selection) {
		selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
			proxyIp := selection.Find("td:nth-child(1)").Text()
			proxyPort := selection.Find("td:nth-child(2)").Text()
			proxyLocation := selection.Find("td:nth-child(4)").Text()
			proxySpeed := selection.Find("td:nth-child(6)").Text()
			ip := new(database.IP)

			ip.ProxyHost = proxyIp
			ip.ProxyPort, _ = strconv.Atoi(proxyPort)
			ip.ProxyType = "https"
			ip.ProxyLocation = proxyLocation
			ip.ProxySpeed, _ = strconv.Atoi(proxySpeed)
			ip.ProxySource = "https://www.zdaye.com"
			ip.CreateTime = util.FormatDateTime()
			ip.UpdateTime = util.FormatDateTime()
			list = append(list, ip)

		})
	})
	log.Info("[Zdaye] fetch done")
	return list
}
