package github

import (
	"IpProxyPool/middleware/database"
	"context"
)

func Yemixzy(ctx context.Context) []*database.IP {
	name := "yemixzy"

	urls := []string{
		"https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/http.txt",
		"https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/socks5.txt",
		"https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/socks4.txt",
		"https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/unchecked.txt",
	}

	return fetchBatch(ctx, name, urls...)
}
