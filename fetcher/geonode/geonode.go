package geonode

import (
	"IpProxyPool/middleware/database"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/youcd/toolkit/log"
)

//nolint:tagliatelle,revive
type item struct {
	Id                 string      `json:"_id"`
	Ip                 string      `json:"ip"`
	AnonymityLevel     string      `json:"anonymityLevel"`
	Asn                string      `json:"asn"`
	City               string      `json:"city"`
	Country            string      `json:"country"`
	CreatedAt          time.Time   `json:"created_at"`
	Google             bool        `json:"google"`
	Isp                string      `json:"isp"`
	LastChecked        int         `json:"lastChecked"`
	Latency            float64     `json:"latency"`
	Org                string      `json:"org"`
	Port               int         `json:"port,string"`
	Protocols          []string    `json:"protocols"`
	Region             interface{} `json:"region"`
	ResponseTime       int         `json:"responseTime"`
	Speed              int         `json:"speed"`
	UpdatedAt          time.Time   `json:"updated_at"`
	WorkingPercent     interface{} `json:"workingPercent"`
	UpTime             float64     `json:"upTime"`
	UpTimeSuccessCount int         `json:"upTimeSuccessCount"`
	UpTimeTryCount     int         `json:"upTimeTryCount"`
}
type resp struct {
	Data []item `json:"data"`
}

func Geonode() []*database.IP {
	const url = "https://proxylist.geonode.com/api/proxy-list?protocols=socks5&limit=500&page=1&sort_by=lastChecked&sort_type=desc"

	// 1. 原生 http 拿字节流，不走进 goquery
	respA, err := http.Get(url)
	if err != nil {
		log.Errorf("geonode http get error: %v", err)
		return nil
	}
	defer respA.Body.Close()

	// 2. 限制最大 2 MB，防止被恶意超大 JSON 打爆
	body, err := io.ReadAll(io.LimitReader(respA.Body, 2<<20))
	if err != nil {
		log.Errorf("geonode read body error: %v", err)
		return nil
	}

	// 3. 直接 JSON 解码
	var respData resp
	if err := json.Unmarshal(body, &respData); err != nil {
		log.Errorf("geonode json error: %v", err)
		return nil
	}

	// 4. 转模型
	list := make([]*database.IP, 0, len(respData.Data))
	for _, d := range respData.Data {
		list = append(list, &database.IP{
			ProxyHost:     d.Ip,
			ProxyPort:     d.Port,
			ProxyType:     d.Protocols[0],
			ProxyLocation: d.City,
			ProxySource:   "https://proxylist.geonode.com",
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		})
	}
	return list
}
