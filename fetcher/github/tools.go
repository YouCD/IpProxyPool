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
		workers    = 100  // 并发数，可根据网络情况调整
		chanBuffer = 1000 // 限制通道容量，避免内存堆积
	)
	lineCh := make(chan string, chanBuffer)
	ipCh := make(chan *database.IP, chanBuffer)

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 5)
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
	for i := 0; i < workers; i++ {
		go func(workerID int) {
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
					}
				case <-ctx.Done():
					return
				}
			}
		}(i)
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
			}
		}
	}()

	// --- 5. 收集结果 ---
	go func() {
		wg.Wait()
		close(ipCh)
	}()

	var list []*database.IP
	for ip := range ipCh {
		list = append(list, ip)
	}

	log.Infof("[%s] fetched %d valid IPs", proxyWeb.Name, len(list))
	return list
}
func fetchBatch(ctx context.Context, name string, urls ...string) []*database.IP {
	list := make([]*database.IP, 0)
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(ctx context.Context, url string) {
			defer wg.Done()
			list = append(list, fetch(ctx, NewProxyWeb(name, url))...)
		}(ctx, url)
	}
	wg.Wait()
	return list
}
