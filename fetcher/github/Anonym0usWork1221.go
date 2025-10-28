package github

import (
	"IpProxyPool/middleware/database"
	"context"
)

func Anonym0usWork1221(ctx context.Context) []*database.IP {
	name := "Anonym0usWork1221"
	urls := []string{
		"https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/http_proxies.txt",
		"https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/https_proxies.txt",
		"https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/socks4_proxies.txt",
		"https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/socks5_proxies.txt",
	}
	return fetchBatch(ctx, name, urls...)
}
