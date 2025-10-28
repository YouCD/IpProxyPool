package github

import (
	"IpProxyPool/middleware/database"
	"context"
)

func R00tee(ctx context.Context) []*database.IP {
	name := "R00tee"
	urls := []string{
		"https://raw.githubusercontent.com/r00tee/Proxy-List/refs/heads/main/Https.txt", // https
		"https://raw.githubusercontent.com/r00tee/Proxy-List/refs/heads/main/Socks4.txt",
		"https://raw.githubusercontent.com/r00tee/Proxy-List/refs/heads/main/Socks5.txt",
	}

	return fetchBatch(ctx, name, urls...)
}
