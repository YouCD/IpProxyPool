package github

import (
	"IpProxyPool/middleware/database"
	"context"
)

func ProxyScraper(ctx context.Context) []*database.IP {
	name := "ProxyScraper"

	urls := []string{
		"https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/http.txt",
		"https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/socks4.txt",
		"https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/socks5.txt",
	}

	return fetchBatch(ctx, name, urls...)
}
