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

func ZloiUser() []*database.IP {
	list := make([]*database.IP, 0)

	name := "zloiUser"
	list = append(list, hideIPMeFetch(NewProxyWeb(name, "https://raw.githubusercontent.com/zloi-user/hideip.me/refs/heads/master/https.txt"))...)

	list = append(list, hideIPMeFetch(NewProxyWeb(name, "https://raw.githubusercontent.com/zloi-user/hideip.me/refs/heads/master/socks4.txt"))...)

	list = append(list, hideIPMeFetch(NewProxyWeb(name, "https://raw.githubusercontent.com/zloi-user/hideip.me/refs/heads/master/socks5.txt"))...)
	return list
}
func hideIPMeFetch(urlStr *ProxyWeb) []*database.IP {
	list := make([]*database.IP, 0, 256)

	doc, _, err := util.Fetch(urlStr.GetFullURL())
	if err != nil {
		log.Errorf("hideip.me fetch %s error: %v", urlStr.GetFullURL(), err)
		return list
	}

	// 流式扫描，避免 document.Text() 一次性分配大字符串
	scanner := bufio.NewScanner(bytes.NewReader([]byte(doc.Text())))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}

		port, _ := strconv.Atoi(parts[1])
		list = append(list, &database.IP{
			ProxyHost:     parts[0],
			ProxyPort:     port,
			ProxyLocation: parts[2],
			ProxySpeed:    100,
			ProxySource:   "https://github.com/zloi-user/hideip.me",
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		})
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("hideip.me scanner error: %v", err)
	}
	return list
}
