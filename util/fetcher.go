package util

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/publicsuffix"
)

var (
	once   sync.Once
	client *http.Client
)

func GetClient() *http.Client {
	once.Do(func() {
		// &cookiejar.Options{PublicSuffixList: publicsuffix.List}，这是为了可以根据域名安全地设置cookies
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			panic(err)
		}
		client = &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second, // 减少超时时间
			Transport: &http.Transport{
				//nolint:gosec
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second, // 减少空闲连接超时时间
			},
		}
	})
	return client
}
func Fetch(ctx context.Context, url string) (*goquery.Document, []byte, error) {
	log.Debugf("Fetch url: %s", url)
	var count int
Retry:
	client := GetClient()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Proxy-Switch-Ip", "yes")
	req.Header.Set("User-Agent", RandomUserAgent())
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "text/html; charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			count++
			if count <= 3 {
				goto Retry
			}
		}
		return nil, nil, fmt.Errorf("fetch url: %s, error: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		count++
		if count <= 3 {
			goto Retry
		}
		return nil, nil, fmt.Errorf("fetch url: %s, status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read body error: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	var reader io.Reader
	if contentType != "" {
		reader, err = charset.NewReader(bytes.NewReader(body), contentType)
		if err != nil {
			reader = bytes.NewReader(body)
		}
	} else {
		reader = bytes.NewReader(body)
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("goquery parse error: %w", err)
	}

	return doc, body, nil
}
