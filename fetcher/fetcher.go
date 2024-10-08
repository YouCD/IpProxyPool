package fetcher

import (
	"IpProxyPool/util"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/publicsuffix"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func Fetch(url string) (*goquery.Document, error) {
	log.Debugf("Fetch url: %s", url)
	var count int
Retry:
	// &cookiejar.Options{PublicSuffixList: publicsuffix.List}，这是为了可以根据域名安全地设置cookies
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}
	//nolint:gosec
	client := &http.Client{
		Jar:     cookieJar,
		Timeout: 100 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	req.Header.Set("Proxy-Switch-Ip", "yes")
	req.Header.Set("User-Agent", util.RandomUserAgent())
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "text/html; charset=UTF-8")

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		count++
		if count < 3 {
			time.Sleep(time.Second * 1)
			goto Retry
		}
		return nil, fmt.Errorf("http get error: %w", err)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("recover get error: %v", err)
		}
	}()

	var newResp io.Reader
	var charsetErr error

	var doc *goquery.Document
	var docErr error

	newResp, charsetErr = charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if charsetErr != nil {
		log.Errorf("charset convert failed: %v", charsetErr)
		return nil, fmt.Errorf("charset convert failed: %w", charsetErr)
	}
	doc, docErr = goquery.NewDocumentFromReader(newResp)
	if docErr != nil {
		log.Errorf("goquery http response body reader error: %v", docErr)
		return nil, fmt.Errorf("goquery http response body reader error: %w", docErr)
	}

	return doc, nil
}
