package github

import (
	"IpProxyPool/middleware/config"
	"IpProxyPool/util"
	"strings"
)

type ProxyWeb struct {
	Name     string
	URL      string
	ProxyURL string
}

func NewProxyWeb(name string, url string) *ProxyWeb {
	p := &ProxyWeb{Name: name, URL: url}
	p.RandomProxy()
	return p
}

func (p *ProxyWeb) RandomProxy() {
	if strings.ContainsAny(p.URL, "raw.githubusercontent.com") {
		userAgentCount := len(config.ServerSetting.GithubProxy)
		randomNum := util.RandInt(0, userAgentCount)
		proxyURL := config.ServerSetting.GithubProxy[randomNum]
		p.ProxyURL = proxyURL
		if strings.HasSuffix(p.ProxyURL, "/") {
			p.ProxyURL = config.ServerSetting.GithubProxy[randomNum]
			return
		}
		p.ProxyURL = config.ServerSetting.GithubProxy[randomNum] + "/"
	}
}

func (p *ProxyWeb) GetFullURL() string {
	if strings.ContainsAny(p.URL, "raw.githubusercontent.com") {
		return p.ProxyURL + p.URL
	}
	return p.URL
}

func (p *ProxyWeb) ChangeProxy() {
	p.RandomProxy()
}
