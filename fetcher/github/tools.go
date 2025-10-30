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
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/youcd/toolkit/log"
)

//nolint:gocognit
func fetch(ctx context.Context, proxyWeb *ProxyWeb) []*database.IP {
	if proxyWeb == nil {
		return nil
	}

	// --- 1. 获取源数据 ---
	_, body, err := util.Fetch(ctx, proxyWeb.GetFullURL())
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		proxyWeb.ChangeProxy()
		log.Errorf("[%s] fetch failed, url: %s, err: %v", proxyWeb.Name, proxyWeb.GetFullURL(), err)
		return nil
	}

	const (
		workers    = 5   // 进一步减少并发数，从20到5
		chanBuffer = 50  // 减少缓冲区大小，避免内存堆积
	)
	lineCh := make(chan string, chanBuffer)
	ipCh := make(chan *database.IP, chanBuffer)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Infof("[%s] ipCh len=%d lineCh len=%d goroutines=%d",
					proxyWeb.Name, len(ipCh), len(lineCh), runtime.NumGoroutine())
			}
		}
	}(ctx)
	var wg sync.WaitGroup
	wg.Add(workers)

	// --- 3. 启动 worker ---
	for range workers {
		go func() {
			defer wg.Done()

			d := net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 2 * time.Second,
			}
			for {
				select {
				case addr, ok := <-lineCh:
					if !ok {
						return
					}

					ipPort := strings.Split(addr, ":")
					if len(ipPort) != 2 {
						continue
					}

					port, err := strconv.Atoi(ipPort[1])
					if err != nil {
						continue
					}

					conn, err := d.DialContext(ctx, "tcp", addr)
					if err != nil {
						continue
					}
					conn.Close()

					select {
					case <-ctx.Done():
						return
					case ipCh <- &database.IP{
						ProxyHost:     ipPort[0],
						ProxyPort:     port,
						ProxyLocation: proxyWeb.Name,
						ProxySpeed:    100,
						ProxySource:   proxyWeb.Name,
						CreateTime:    time.Now(),
						UpdateTime:    time.Now(),
					}:
					case <-time.After(2 * time.Second):
						log.Warnf("send timeout, channel likely full: %s, len: %d", "ipCh", len(ipCh))
						return // 添加超时返回机制
					}
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					// 防止长时间阻塞
					return
				}
			}
		}()
	}

	// --- 4. 边读边派发任务 ---
	go func() {
		defer close(lineCh)

		scanner := bufio.NewScanner(bytes.NewReader(body))
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			addr := strings.TrimSpace(scanner.Text())
			if addr == "" {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case lineCh <- addr:
			case <-time.After(2 * time.Second):
				log.Warnf("send timeout, channel likely full: %s, len: %d", "lineCh", len(lineCh))
				return // 添加超时返回机制
			}
		}
	}()

	// --- 5. 收集结果 ---
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(ipCh)
		close(done)
	}()

	list := make([]*database.IP, 0, len(ipCh))
	for {
		select {
		case ip, ok := <-ipCh:
			if !ok {
				// Channel closed, all IPs collected
				goto finish
			}
			list = append(list, ip)
		case <-done:
			// Workers finished, collect remaining IPs
			for {
				select {
				case ip, ok := <-ipCh:
					if !ok {
						goto finish
					}
					list = append(list, ip)
				default:
					goto finish
				}
			}
		case <-ctx.Done():
			goto finish
		}
	}

finish:
	log.Debugf("[%s] fetched %d valid IPs", proxyWeb.Name, len(list))
	return list
}

func fetchBatch(ctx context.Context, name string, urls ...string) []*database.IP {
	list := make([]*database.IP, 0)
	var wg sync.WaitGroup
	
	// 限制并发获取URL的数量
	const maxConcurrentFetches = 3
	semaphore := make(chan struct{}, maxConcurrentFetches)

	for _, url := range urls {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量
		go func(ctx context.Context, url string) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量
			list = append(list, fetch(ctx, NewProxyWeb(name, url))...)
		}(ctx, url)
	}
	wg.Wait()
	return list
}