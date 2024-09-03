package geonode

import (
	"IpProxyPool/fetcher"
	"IpProxyPool/middleware/database"
	"IpProxyPool/util"
	"encoding/json"
	"github.com/youcd/toolkit/log"
	"time"
)

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
	log.Info("[geonode] fetch start")
	list := make([]*database.IP, 0)
	document, err := fetcher.Fetch("https://proxylist.geonode.com/api/proxy-list?protocols=socks5&limit=500&page=1&sort_by=lastChecked&sort_type=desc")
	if err != nil {
		log.Errorf("document geonode error:%s", err)
		return list
	}
	var respData resp
	if err := json.Unmarshal([]byte(document.Text()), &respData); err != nil {
		log.Error(err)
		return list
	}
	for _, datum := range respData.Data {
		ip := &database.IP{
			ProxyHost:     datum.Ip,
			ProxyPort:     datum.Port,
			ProxyType:     datum.Protocols[0],
			ProxyLocation: datum.City,
			ProxySource:   "https://proxylist.geonode.com",
			CreateTime:    util.FormatDateTime(),
			UpdateTime:    util.FormatDateTime(),
		}
		list = append(list, ip)
	}
	log.Info("[geonode] fetch done")
	return list
}
