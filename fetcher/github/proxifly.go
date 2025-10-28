package github

import (
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"bufio"
	"bytes"
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/youcd/toolkit/log"
)

func FreeProxyList(ctx context.Context) []*database.IP {
	list := make([]*database.IP, 0)
	name := "FreeProxyList"
	var wg sync.WaitGroup
	urls := []string{
		"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks5/data.txt",
		"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks4/data.txt",
		"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/http/data.txt",
		"https://raw.githubusercontent.com/proxifly/free-proxy-list/refs/heads/main/proxies/countries/US/data.txt",
	}

	for _, u := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			list = append(list, freeProxyListFetch(ctx, NewProxyWeb(name, u))...)
		}(u)
	}
	wg.Wait()

	return list
}
func freeProxyListFetch(ctx context.Context, urlStr *ProxyWeb) []*database.IP {
	list := make([]*database.IP, 0, 256)

	doc, _, err := util.Fetch(ctx, urlStr.GetFullURL())
	if err != nil {
		log.Errorf("%s fetch error: %v", urlStr.Name, err)
		return list
	}

	// 流式扫描，避免 document.Text() 一次性分配大字符串
	scanner := bufio.NewScanner(bytes.NewReader([]byte(doc.Text())))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parse, err := url.Parse(line)
		if err != nil {
			log.Error("free-proxy-list parse error: %v", err)
			continue
		}

		list = append(list, &database.IP{
			ProxyHost:     parse.Host,
			ProxyPort:     cast.ToInt(parse.Port()),
			ProxyLocation: "free-proxy-list",
			ProxyType:     parse.Scheme,
			ProxySpeed:    100,
			ProxySource:   "https://github.com/proxifly/free-proxy-list",
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		})
	}
	err = scanner.Err()
	if err != nil {
		log.Errorf("free-proxy-list scanner error: %v", err)
	}
	return list
}
