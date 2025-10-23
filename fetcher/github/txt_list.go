package github

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/youcd/toolkit/log"
)

func fetch(proxyWeb *ProxyWeb) []*database.IP {
	if proxyWeb == nil {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			log.Warnf("[%s] fetch panic: %v", proxyWeb.Name, r)
		}
	}()

	var count int
Retry:
	document, err := fetcher.Fetch(proxyWeb.GetFullURL())
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		proxyWeb.ChangeProxy()
		count++
		if count < 3 {
			log.Errorf("[%s] ChangeProxy: %s ", proxyWeb.Name, proxyWeb.ProxyURL)
			goto Retry
		}
		log.Errorf("[%s] fetch failed, url: %s, err: %v", proxyWeb.Name, proxyWeb.GetFullURL(), err)
		return nil
	}

	// 使用 bufio.Scanner 逐行读取
	scanner := bufio.NewScanner(strings.NewReader(document.Text()))
	ch := make(chan *database.IP, 1024)
	var wg sync.WaitGroup

	// 控制最大并发数（防止 goroutine 爆炸）
	const maxConcurrent = 50
	sem := make(chan struct{}, maxConcurrent)

	for scanner.Scan() {
		address := strings.TrimSpace(scanner.Text())
		if address == "" {
			continue
		}

		wg.Add(1)
		sem <- struct{}{} // 占位，限制并发

		go func(ipPort string) {
			defer func() {
				<-sem     // 释放并发槽
				wg.Done() // 结束信号
			}()

			conn, err := net.DialTimeout("tcp", ipPort, 3*time.Second)
			if err != nil {
				return
			}
			_ = conn.Close()

			ipPortObj := strings.Split(ipPort, ":")
			if len(ipPortObj) != 2 {
				return
			}
			port, err := strconv.Atoi(ipPortObj[1])
			if err != nil {
				return
			}

			ip := &database.IP{
				ProxyHost:     ipPortObj[0],
				ProxyPort:     port,
				ProxyLocation: proxyWeb.Name,
				ProxySpeed:    100,
				ProxySource:   proxyWeb.Name,
				CreateTime:    time.Now(),
				UpdateTime:    time.Now(),
			}
			ch <- ip
		}(address)
	}

	wg.Wait()
	close(ch)

	// 汇总结果
	list := make([]*database.IP, 0, len(ch))
	for ip := range ch {
		list = append(list, ip)
	}

	return list
}
