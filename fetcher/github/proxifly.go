package github

import (
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"bufio"
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/youcd/toolkit/log"
)

func FreeProxyList() []*database.IP {
	list := make([]*database.IP, 0)
	name := "FreeProxyList"
	socks5Url := NewProxyWeb(name, "https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks5/data.txt")
	list = append(list, freeProxyListFetch(socks5Url, "socks5://")...)

	socks4Url := NewProxyWeb(name, "https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks4/data.txt")
	list = append(list, freeProxyListFetch(socks4Url, "socks4://")...)

	httpURL := NewProxyWeb(name, "https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/http/data.txt")
	list = append(list, freeProxyListFetch(httpURL, "http://")...)
	return list
}
func freeProxyListFetch(urlStr *ProxyWeb, replaceStr string) []*database.IP {
	list := make([]*database.IP, 0, 256)

	doc, _, err := util.Fetch(urlStr.GetFullURL())
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
		line = strings.ReplaceAll(line, replaceStr, "")
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		port, _ := strconv.Atoi(parts[1])
		list = append(list, &database.IP{
			ProxyHost:     parts[0],
			ProxyPort:     port,
			ProxyLocation: "free-proxy-list",
			ProxySpeed:    100,
			ProxySource:   "https://github.com/proxifly/free-proxy-list",
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		})
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("free-proxy-list scanner error: %v", err)
	}
	return list
}
