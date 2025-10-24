package github

import (
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"bufio"
	"bytes"
	"context"
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

	var count int
Retry:
	_, body, err := util.Fetch(proxyWeb.GetFullURL())
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		proxyWeb.ChangeProxy()
		count++
		if count < 3 {
			log.Errorf("[%s] ChangeProxy: %s", proxyWeb.Name, proxyWeb.ProxyURL)
			goto Retry
		}
		log.Errorf("[%s] fetch failed, url: %s, err: %v", proxyWeb.Name, proxyWeb.GetFullURL(), err)
		return nil
	}

	// 1. 收集有效行
	scanner := bufio.NewScanner(bytes.NewReader(body))
	var lines []string
	for scanner.Scan() {
		if addr := strings.TrimSpace(scanner.Text()); addr != "" {
			lines = append(lines, addr)
		}
	}

	if len(lines) == 0 {
		return nil
	}

	// 2. 使用context控制超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const workers = 50
	lineCh := make(chan string, len(lines))
	ipCh := make(chan *database.IP, len(lines))
	var wg sync.WaitGroup
	wg.Add(workers)

	// 启动worker
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case addr, ok := <-lineCh:
					if !ok {
						return // 通道关闭，退出
					}
					ipPort := strings.Split(addr, ":")
					if len(ipPort) != 2 {
						continue
					}
					port, err := strconv.Atoi(ipPort[1])
					if err != nil {
						continue
					}

					// 使用带context的Dial
					var d net.Dialer
					conn, err := d.DialContext(ctx, "tcp", addr)
					if err != nil {
						continue
					}
					conn.Close()

					select {
					case ipCh <- &database.IP{
						ProxyHost:     ipPort[0],
						ProxyPort:     port,
						ProxyLocation: proxyWeb.Name,
						ProxySpeed:    100,
						ProxySource:   proxyWeb.Name,
						CreateTime:    time.Now(),
						UpdateTime:    time.Now(),
					}:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return // 超时退出
				}
			}
		}(i)
	}

	// 3. 分发任务（带超时控制）
	go func() {
		defer close(lineCh)
		for _, l := range lines {
			select {
			case lineCh <- l:
			case <-ctx.Done():
				return
			}
		}
	}()

	// 4. 等待worker完成（带超时控制）
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// 等待worker完成或超时
	select {
	case <-done:
		// worker正常完成
	case <-ctx.Done():
		log.Warnf("[%s] fetch timeout, url: %s", proxyWeb.Name, proxyWeb.GetFullURL())
	}

	close(ipCh)

	// 5. 汇总结果
	list := make([]*database.IP, 0, len(ipCh))
	for ip := range ipCh {
		list = append(list, ip)
	}

	log.Infof("[%s] fetched %d valid IPs from %d lines", proxyWeb.Name, len(list), len(lines))
	return list
}
