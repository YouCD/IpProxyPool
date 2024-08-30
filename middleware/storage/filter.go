package storage

import (
	"IpProxyPool/middleware/database"
	"IpProxyPool/util/randomutil"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/youcd/toolkit/log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	httpsHttpBin = "https://www.google.com"
	httpHttpBin  = "http://httpbin.org/get"
)

// CheckProxy .
func CheckProxy(ip *database.IP) {
	if checkIp, ok := CheckIp(ip); ok {
		database.SaveIp(checkIp)
		log.Debug("proxy is good           ", checkIp)
	}
}

// CheckIp
//
//	@Description: 检测代理IP是否可用
//	@param ip
//	@return *database.IP
//	@return bool
func CheckIp(ip *database.IP) (*database.IP, bool) {
	d, err := checkIP(ip)
	if err != nil {
		log.Debug(err)
		return nil, false
	}

	return d, true
}

func checkIP(ip *database.IP) (*database.IP, error) {
	if ip == nil {
		return nil, errors.New("proxy is empty")
	}
	// 解析为 https 代理
	if isHttpsProxy(ip) {
		ip.ProxyType = "https"
		return ip, nil
	}
	// 解析为 http 代理
	if isHttpProxy(ip) {
		ip.ProxyType = "http"
		return ip, nil
	}

	// 解析为 socks5 代理
	if isSocks5Proxy(ip) {
		ip.ProxyType = "socks5"
		return ip, nil
	}
	// 解析为 tcp 代理
	if isTcpProxy(ip) {
		ip.ProxyType = "tcp"
		return ip, nil
	}

	return nil, errors.New("proxy is not available")
}

func isSocks5Proxy(ip *database.IP) bool {
	return requestHTTPBIN(ip, httpsHttpBin, "socks5")
}

func isTcpProxy(ip *database.IP) bool {
	return requestHTTPBIN(ip, httpsHttpBin, "tcp")
}

func isHttpProxy(ip *database.IP) bool {
	return requestHTTPBIN(ip, httpHttpBin, "http")
}

func isHttpsProxy(ip *database.IP) bool {
	return requestHTTPBIN(ip, httpsHttpBin, "https")
}
func requestHTTPBIN(ip *database.IP, testUrl string, scheme string) bool {
	address := fmt.Sprintf("%s:%d", ip.ProxyHost, ip.ProxyPort)
	proxy, err := url.Parse(strings.ToLower(fmt.Sprintf("%s://%s", scheme, address)))
	if err != nil {
		return false
	}

	dialer := &net.Dialer{
		// 限制创建一个TCP连接使用的时间（如果需要一个新的链接）
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	// 设置网络传输
	netTransport := &http.Transport{
		DialContext:           dialer.DialContext,
		Proxy:                 http.ProxyURL(proxy),
		DisableKeepAlives:     true,
		MaxConnsPerHost:       20,
		MaxIdleConns:          20,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       20 * time.Second,
		ResponseHeaderTimeout: time.Second * time.Duration(10),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	// 创建连接客户端
	httpClient := &http.Client{
		Transport: netTransport,
	}

	begin := time.Now() //判断代理访问时间

	// 使用代理IP访问测试地址
	res, err := httpClient.Get(testUrl)
	if err != nil {
		log.Debugf("testIp: %s, testUrl: %s: error: %v", address, testUrl, err.Error())
		return false
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		// 判断是否成功访问，如果成功访问StatusCode应该为200
		speed := time.Now().Sub(begin).Nanoseconds() / 1000 / 1000 //ms
		ip.ProxySpeed = int(speed)
		database.UpdateIp(ip)
		return true
	}
	return false
}

// CheckProxyDB to check the ip in DB
func CheckProxyDB() {
	record := database.CountIp()
	log.Infof("Before check, DB has: %d records.", record)
	ips := database.GetAllIp()
	var wg sync.WaitGroup
	for _, v := range ips {
		wg.Add(1)
		go func(ip *database.IP) {
			defer wg.Done()
			newIP, ok := CheckIp(ip)
			if !ok {
				database.DeleteIp(ip)
			} else {
				database.UpdateIp(newIP)
			}
		}(v)
	}
	wg.Wait()
	record = database.CountIp()
	log.Infof("After check, DB has: %d records.", record)
}

// RandomProxy .
func RandomProxy() (ip *database.IP) {
	ips := database.GetAllIp()
	ipCount := len(ips)
	if ipCount == 0 {
		log.Warnf("RandomProxy random count: %d\n", ipCount)
		return nil
	}
	randomNum := randomutil.RandInt(0, ipCount)
	return ips[randomNum]
}

// RandomByProxyType .
func RandomByProxyType(proxyType string) (ip database.IP) {
	//  如果是 https 优先返回 tcp的代理类型
	switch proxyType {
	case "https":
		var ips []database.IP
		// tcp
		ipsForTcp, err := database.GetIpByProxyType("tcp")
		if err == nil {
			ips = append(ips, ipsForTcp...)
		}
		log.Debugf("proxy_type: tcp, count: %d ", len(ipsForTcp))

		// socks5
		ipsForSocks5, err := database.GetIpByProxyType("socks5")
		if err == nil {
			ips = append(ips, ipsForSocks5...)
		}
		log.Debugf("proxy_type: socks5, count: %d ", len(ipsForSocks5))

		// https
		ipsForHttps, err := database.GetIpByProxyType("https")
		if err == nil {
			ips = append(ips, ipsForHttps...)
		}
		log.Debugf("proxy_type: https, count: %d ", len(ipsForHttps))

		// socks4
		ipsForSocks4, err := database.GetIpByProxyType("socks4")
		if err == nil {
			ips = append(ips, ipsForSocks4...)
		}
		log.Debugf("proxy_type: socks4, count: %d ", len(ipsForSocks4))

		ipCount := len(ips)
		if ipCount == 0 {
			//  如果没有 tcp 类型的代理，就返回 https 类型的代理
			return randomByProxyType(proxyType)
		}
		randomNum := randomutil.RandInt(0, ipCount)
		return ips[randomNum]
	case "http":
		// http
		ipsForHttp, err := database.GetIpByProxyType("http")
		if err != nil {
			return randomByProxyType(proxyType)
		}

		ipCount := len(ipsForHttp)
		if ipCount == 0 {
			//  如果没有 tcp 类型的代理，就返回 https 类型的代理
			return randomByProxyType(proxyType)
		}

		randomNum := randomutil.RandInt(0, ipCount)
		log.Debugf("proxy_type: http, count: %d ", len(ipsForHttp))
		return ipsForHttp[randomNum]
	default:
		return randomByProxyType(proxyType)
	}

}
func randomByProxyType(proxyType string) (ip database.IP) {
	ips, err := database.GetIpByProxyType(proxyType)
	if err != nil {
		log.Warn(err.Error())
		return database.IP{}
	}
	ipCount := len(ips)
	if ipCount == 0 {
		log.Warnf("RandomByProxyType random count: %d\n", ipCount)
		return database.IP{}
	}
	randomNum := randomutil.RandInt(0, ipCount)
	return ips[randomNum]
}
