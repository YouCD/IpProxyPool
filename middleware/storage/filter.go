package storage

import (
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"context"
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
	httpsHTTPBin = "https://www.google.com"
	httpHTTPBin  = "http://httpbin.org/get"
)

var (
	ErrNotAvailable = errors.New("proxy is not available")
	ErrProxyEmpty   = errors.New("proxy is empty")
)

// CheckProxy .
func CheckProxy(ip *database.IP) {
	now := time.Now()
	if ip == nil {
		log.Error("CheckProxy empty ip")
		return
	}
	var flag bool
	if item, ok := CheckIP(ip); ok {
		database.SaveIP(item)
		flag = true
	}
	log.Infof("Checking proxy: %s://%s:%d, isOK: %t, Duration: %s", ip.ProxyType, ip.ProxyHost, ip.ProxyPort, flag, time.Since(now))
}

// CheckIP
//
//	@Description: 检测代理IP是否可用
//	@param ip
//	@return *database.IP
//	@return bool
func CheckIP(ip *database.IP) (*database.IP, bool) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 60*time.Second)
	defer func() {
		//nolint:gosimple
		select {
		case <-ctx.Done():
			time.Sleep(5 * time.Second)
			cancelFunc()
		}
	}()
	d := checkIP(ctx, ip)
	if d == nil {
		return ip, false
	}
	return d, true
}
func checkIP(ctx context.Context, ip *database.IP) *database.IP {
	if ip == nil {
		return nil
	}
	if _, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip.ProxyHost, ip.ProxyPort), 3*time.Second); err != nil {
		return nil
	}
	resultChan := make(chan *database.IP, 1)
	var wg sync.WaitGroup
	ip.ProxyType = strings.ToLower(ip.ProxyType)

	// Helper function to send results to channels
	sendResult := func(ctx context.Context, proxyType string, proxyCheckFunc func(*database.IP) bool) {
		defer wg.Done() // Decrease WaitGroup counter when done
		select {
		case <-ctx.Done():
			return
		default:
			if proxyCheckFunc(ip) {
				ip.ProxyType = proxyType
				select {
				case resultChan <- ip:
				default: // If resultChan already has a result, skip
				}
			}
		}
	}

	wg.Add(5)
	go sendResult(ctx, "https", isHTTPSProxy)
	go sendResult(ctx, "http", isHTTPProxy)
	go sendResult(ctx, "socks5", isSocks5Proxy)
	go sendResult(ctx, "socks4", isSocks4Proxy)
	go sendResult(ctx, "tcp", isTCPProxy)

	go func() {
		wg.Wait()         // Wait for all goroutines to finish
		close(resultChan) // Close the channel after all goroutines finish
	}()
	return <-resultChan
}
func isSocks5Proxy(ip *database.IP) bool {
	return requestHTTPBIN(ip, httpsHTTPBin, "socks5")
}

func isTCPProxy(ip *database.IP) bool {
	return requestHTTPBIN(ip, httpsHTTPBin, "tcp")
}
func isSocks4Proxy(ip *database.IP) bool {
	return requestHTTPBIN(ip, httpsHTTPBin, "socks4")
}

func isHTTPProxy(ip *database.IP) bool {
	return requestHTTPBIN(ip, httpHTTPBin, "http")
}

func isHTTPSProxy(ip *database.IP) bool {
	return requestHTTPBIN(ip, httpsHTTPBin, "https")
}
func requestHTTPBIN(ip *database.IP, testURL string, scheme string) bool {
	address := fmt.Sprintf("%s:%d", ip.ProxyHost, ip.ProxyPort)
	rawURL := strings.ToLower(fmt.Sprintf("%s://%s", scheme, address))
	proxy, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	dialer := &net.Dialer{
		// 限制创建一个TCP连接使用的时间（如果需要一个新的链接）
		Timeout:   10 * time.Second,
		KeepAlive: 10 * time.Second,
	}
	// 设置网络传输
	//nolint:gosec
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
		// DialTLSContext: func(_ context.Context, network, addr string) (net.Conn, error) {
		// 	return tls.Dial(network, addr, &tls.Config{InsecureSkipVerify: true})
		// },
	}

	// 创建连接客户端
	httpClient := &http.Client{
		Transport: netTransport,
	}

	begin := time.Now() // 判断代理访问时间

	// 使用代理IP访问测试地址
	//nolint:noctx
	res, err := httpClient.Get(testURL)
	if err != nil {
		log.Debugf("proxy: %s, error: %s", rawURL, err)
		return false
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		// 判断是否成功访问，如果成功访问StatusCode应该为200
		speed := time.Since(begin).Nanoseconds() / 1000 / 1000 // ms
		ip.ProxySpeed = int(speed)
		database.UpdateIP(ip)
		return true
	}
	return false
}

// CheckProxyDB to check the ip in DB
func CheckProxyDB() {
	start := time.Now()
	beforeRecord := database.CountIP()
	ips := database.GetAllIP()
	var wg sync.WaitGroup
	for _, v := range ips {
		wg.Add(1)
		go func(ip *database.IP) {
			defer wg.Done()
			newIP, ok := CheckIP(ip)
			if !ok {
				log.Warnf("CheckProxyDB proxy: %s, error: %s", ip.ProxyHost, ErrNotAvailable)
				database.DeleteIP(ip)
			} else {
				database.UpdateIP(newIP)
			}
		}(v)
	}
	wg.Wait()
	afterRecord := database.CountIP()
	log.Infof("Before cout: %d, After cout:%d, Duration: %s", beforeRecord, afterRecord, time.Since(start))
}

// RandomProxy .
func RandomProxy() (ip *database.IP) {
	ips := database.GetAllIP()
	ipCount := len(ips)
	if ipCount == 0 {
		log.Warnf("RandomProxy random count: %d\n", ipCount)
		return nil
	}
	randomNum := util.RandInt(0, ipCount)
	return ips[randomNum]
}

// RandomByProxyType .
func RandomByProxyType(proxyType string) (ip database.IP) {
	//  如果是 https 优先返回 tcp的代理类型
	switch proxyType {
	case "https":
		var ips []database.IP
		// tcp
		ipsForTCP, err := database.GetIPByProxyType("tcp")
		if err == nil {
			ips = append(ips, ipsForTCP...)
		}
		log.Debugf("proxy_type: tcp, count: %d ", len(ipsForTCP))

		// socks5
		ipsForSocks5, err := database.GetIPByProxyType("socks5")
		if err == nil {
			ips = append(ips, ipsForSocks5...)
		}
		log.Debugf("proxy_type: socks5, count: %d ", len(ipsForSocks5))

		// https
		ipsForHTTPS, err := database.GetIPByProxyType("https")
		if err == nil {
			ips = append(ips, ipsForHTTPS...)
		}
		log.Debugf("proxy_type: https, count: %d ", len(ipsForHTTPS))

		// socks4
		ipsForSocks4, err := database.GetIPByProxyType("socks4")
		if err == nil {
			ips = append(ips, ipsForSocks4...)
		}
		log.Debugf("proxy_type: socks4, count: %d ", len(ipsForSocks4))

		ipCount := len(ips)
		if ipCount == 0 {
			//  如果没有 tcp 类型的代理，就返回 https 类型的代理
			return randomByProxyType(proxyType)
		}
		randomNum := util.RandInt(0, ipCount)
		return ips[randomNum]
	case "http":
		// http
		ipsForHTTP, err := database.GetIPByProxyType("http")
		if err != nil {
			return randomByProxyType(proxyType)
		}

		ipCount := len(ipsForHTTP)
		if ipCount == 0 {
			//  如果没有 tcp 类型的代理，就返回 https 类型的代理
			return randomByProxyType(proxyType)
		}

		randomNum := util.RandInt(0, ipCount)
		log.Debugf("proxy_type: http, count: %d ", len(ipsForHTTP))
		return ipsForHTTP[randomNum]
	default:
		return randomByProxyType(proxyType)
	}
}
func randomByProxyType(proxyType string) (ip database.IP) {
	ips, err := database.GetIPByProxyType(proxyType)
	if err != nil {
		log.Warn(err)
		return database.IP{}
	}
	ipCount := len(ips)
	if ipCount == 0 {
		log.Warnf("RandomByProxyType random count: %d\n", ipCount)
		return database.IP{}
	}
	randomNum := util.RandInt(0, ipCount)
	return ips[randomNum]
}
