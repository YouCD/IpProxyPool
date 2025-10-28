package github

import (
	"IpProxyPool/middleware/database"
	"context"
)

func Zaeem20(ctx context.Context) []*database.IP {
	name := "Zaeem20"
	urls := []string{
		"https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/http.txt",
		"https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/https.txt",
		"https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/socks4.txt",
	}

	return fetchBatch(ctx, name, urls...)
}
