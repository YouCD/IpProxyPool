package github

import (
	"IpProxyPool/middleware/database"
	"context"
)

func Zenjahid(ctx context.Context) []*database.IP {
	name := "zenjahid"
	urls := []string{
		"https://raw.githubusercontent.com/zenjahid/FreeProxy4u/main/http.txt",
		"https://raw.githubusercontent.com/zenjahid/FreeProxy4u/main/socks4.txt",
		"https://raw.githubusercontent.com/zenjahid/FreeProxy4u/main/socks5.txt",
	}

	return fetchBatch(ctx, name, urls...)
}
