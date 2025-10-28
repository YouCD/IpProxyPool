package github

import (
	"IpProxyPool/middleware/database"
	"context"
)

func Vakhov(ctx context.Context) []*database.IP {
	name := "Vakhov"

	urls := []string{
		"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/https.txt",
		"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/socks5.txt",
		"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/socks4.txt",
		"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/http.txt",
	}

	return fetchBatch(ctx, name, urls...)
}
