package storage

import (
	"IpProxyPool/middleware/database"
	"IpProxyPool/util/randomutil"
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
)

const (
	httpsHttpBin = "https://httpbin.org/get"
	httpHttpBin  = "http://httpbin.org/get"
)

// CheckProxy .
func CheckProxy(ip *database.IP) {
	checkIp, ok := CheckIp(ip)
	if ok {
		database.SaveIp(checkIp)
	}
	logger.Println("check proxy done           ", checkIp)
}

//
// CheckIp
//  @Description: 检测代理IP是否可用
//  @param ip
//  @return *database.IP
//  @return bool
//
func CheckIp(ip *database.IP) (*database.IP, bool) {
	// 检测代理iP访问地址
	var testUrl string
	if ip == nil {
		logger.Error("ip is nil")
		return nil, false
	}
	if ip.ProxyType == "http" {
		testUrl = httpHttpBin
	} else {
		testUrl = httpsHttpBin
	}

	for {
		dialer, Proxy, newIP, err := CreateDialerAndProxy(ip)
		if err != nil {
			logger.Error(err)
			return newIP, false
		}
		testIp := fmt.Sprintf("%s://%s:%d", newIP.ProxyType, newIP.ProxyHost, newIP.ProxyPort)
		// 设置网络传输
		netTransport := &http.Transport{
			DialContext:           dialer.DialContext,
			Proxy:                 Proxy,
			DisableKeepAlives:     true,
			MaxConnsPerHost:       20,
			MaxIdleConns:          20,
			MaxIdleConnsPerHost:   20,
			IdleConnTimeout:       20 * time.Second,
			ResponseHeaderTimeout: time.Second * time.Duration(10),
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		}
		_ = http2.ConfigureTransport(netTransport)

		// 创建连接客户端
		httpClient := &http.Client{
			Transport: netTransport,
		}
		begin := time.Now() //判断代理访问时间

		// 使用代理IP访问测试地址
		res, err := httpClient.Get(testUrl)
		if err != nil {
			logger.Debugf("testIp: %s, testUrl: %s: error msg: %v", testIp, testUrl, err.Error())
		} else {
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				// 判断是否成功访问，如果成功访问StatusCode应该为200
				speed := time.Now().Sub(begin).Nanoseconds() / 1000 / 1000 //ms
				ip.ProxySpeed = int(speed)
				database.UpdateIp(ip)
				return ip, true
			}
		}

		// 切换协议类型
		if ip.ProxyType == "https" {
			ip.ProxyType = "tcp"
		} else {
			break
		}
	}

	return ip, false
}

func CreateDialerAndProxy(ip *database.IP) (*net.Dialer, func(*http.Request) (*url.URL, error), *database.IP, error) {
	if ip == nil {
		return nil, nil, ip, errors.New("proxy is empty")
	}

	ip.ProxyType = strings.ToLower(ip.ProxyType)

	dialer := &net.Dialer{
		// 限制创建一个TCP连接使用的时间（如果需要一个新的链接）
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	switch {
	case strings.Contains(ip.ProxyType, "http"), strings.Contains(ip.ProxyType, "https"):
		// 解析代理地址
		proxy, err := url.Parse(strings.ToLower(fmt.Sprintf("%s://%s:%d", ip.ProxyType, ip.ProxyHost, ip.ProxyPort)))
		if err != nil {
			return nil, nil, ip, err
		}
		return dialer, http.ProxyURL(proxy), ip, nil
	case strings.Contains(ip.ProxyType, "tcp"):
		conn, err := dialer.Dial("tcp", ip.ProxyHost+":"+strconv.Itoa(ip.ProxyPort))
		if err != nil {
			logger.Error(err)
			return nil, nil, ip, err
		}

		defer conn.Close()
		//  变更代理类型为 tcp
		ip.ProxyType = "tcp"
		return dialer, nil, ip, err
	default:
		return nil, nil, ip, nil
	}
}

// CheckProxyDB to check the ip in DB
func CheckProxyDB() {
	record := database.CountIp()
	logger.Infof("Before check, DB has: %d records.", record)
	ips := database.GetAllIp()

	var wg sync.WaitGroup
	for _, v := range ips {
		wg.Add(1)
		go func(ip database.IP) {
			newIP, ok := CheckIp(&ip)
			if !ok {
				database.DeleteIp(&ip)
			}

			database.UpdateIp(newIP)
			wg.Done()
		}(v)
	}
	wg.Wait()
	record = database.CountIp()
	logger.Infof("After check, DB has: %d records.", record)
}

// AllProxy .
func AllProxy() []database.IP {
	ips := database.GetAllIp()
	ipCount := len(ips)
	if ipCount == 0 {
		logger.Warnf("RandomProxy random count: %d\n", ipCount)
		return []database.IP{}
	}
	return ips
}

// RandomProxy .
func RandomProxy() (ip database.IP) {
	ips := database.GetAllIp()
	ipCount := len(ips)
	if ipCount == 0 {
		logger.Warnf("RandomProxy random count: %d\n", ipCount)
		return database.IP{}
	}
	randomNum := randomutil.RandInt(0, ipCount)
	return ips[randomNum]
}

// RandomByProxyType .
func RandomByProxyType(proxyType string) (ip database.IP) {
	//  如果是 https 优先返回 tcp的代理类型
	switch proxyType {
	case "https":
		ips, err := database.GetIpByProxyType("tcp")
		if err != nil {
			logger.Warn(err.Error())
			return database.IP{}
		}
		ipCount := len(ips)
		if ipCount == 0 {
			//  如果没有 tcp 类型的代理，就返回 https 类型的代理
			return randomByProxyType(proxyType)
		}
		randomNum := randomutil.RandInt(0, ipCount)
		return ips[randomNum]
	default:
		return randomByProxyType(proxyType)
	}

}
func randomByProxyType(proxyType string) (ip database.IP) {
	ips, err := database.GetIpByProxyType(proxyType)
	if err != nil {
		logger.Warn(err.Error())
		return database.IP{}
	}
	ipCount := len(ips)
	if ipCount == 0 {
		logger.Warnf("RandomByProxyType random count: %d\n", ipCount)
		return database.IP{}
	}
	randomNum := randomutil.RandInt(0, ipCount)
	return ips[randomNum]
}
