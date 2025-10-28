package github

import (
	"IpProxyPool/middleware/database"
	"context"
)

func TheSpeedX(ctx context.Context) []*database.IP {
	name := "TheSpeedX"

	urls := []string{
		"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt",
		"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks4.txt",
		"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks5.txt",
	}

	return fetchBatch(ctx, name, urls...)
}
