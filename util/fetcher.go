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
			Timeout: 100 * time.Second,
			Transport: &http.Transport{
				//nolint:gosec
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
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
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		count++
		if count < 3 {
			time.Sleep(time.Second * 1)
			goto Retry
		}
		return nil, nil, fmt.Errorf("http get error: %w", err)
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
		if errors.Is(charsetErr, io.EOF) {
			return nil, nil, fmt.Errorf("charset error: %w", charsetErr)
		}
		log.Errorf("charset convert failed: %v", charsetErr)
		return nil, nil, fmt.Errorf("charset convert failed: %w", charsetErr)
	}
	var bufer bytes.Buffer
	_, _ = io.Copy(&bufer, newResp)
	// 先获取字节数据
	bt := bufer.Bytes()

	doc, docErr = goquery.NewDocumentFromReader(bytes.NewReader(bt))
	if docErr != nil {
		log.Errorf("goquery http response body reader error: %v", docErr)
		return nil, nil, fmt.Errorf("goquery http response body reader error: %w", docErr)
	}

	return doc, bt, nil
}
