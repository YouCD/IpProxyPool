package github

import (
	"IpProxyPool/middleware/database"
	"context"
)

func OpenProxyList(ctx context.Context) []*database.IP {
	name := "OpenProxyList"
	urls := []string{
		"https://api.openproxylist.xyz/https.txt",
		"https://api.openproxylist.xyz/http.txt",
		"https://api.openproxylist.xyz/socks5.txt",
		"https://api.openproxylist.xyz/socks4.txt",
	}
	return fetchBatch(ctx, name, urls...)
}
