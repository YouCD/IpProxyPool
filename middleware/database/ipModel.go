package database

import (
	"github.com/youcd/toolkit/log"
	"strings"
	"time"
)

// IP struct
//
//nolint:revive
type IP struct {
	ProxyId       int64     `gorm:"primary_key; auto_increment; not null" json:"-"`
	ProxyHost     string    `gorm:"type:varchar(255); not null; uniqueIndex:UNIQUE_HOST_PORT" json:"proxyHost"`
	ProxyPort     int       `gorm:"type:int(11); not null; uniqueIndex:UNIQUE_HOST_PORT" json:"proxyPort"`
	ProxyType     string    `gorm:"type:varchar(64); not null" json:"proxyType"`
	ProxyLocation string    `gorm:"type:varchar(255); default null" json:"proxyLocation"`
	ProxySpeed    int       `gorm:"type:int(20); not null; default 0" json:"proxySpeed"`
	ProxySource   string    `gorm:"type:varchar(64); not null;" json:"proxySource"`
	CreateTime    time.Time `gorm:"type:time; not null" json:"-"`
	UpdateTime    time.Time `gorm:"type:time; default ''" json:"updateTime"`
}

func (i *IP) TableName() string {
	return "ip"
}

// SaveIP 保存数据到数据库
func SaveIP(ip *IP) {
	db := GetDB().Begin()
	ip.ProxyType = strings.ToLower(ip.ProxyType)
	ipModel := GetIPByProxyHost(ip.ProxyHost)
	if ipModel.ProxyHost == "" {
		err := db.Model(&IP{}).Create(ip).Error
		if err != nil {
			log.Errorf("save ip: %s, error msg: %v", ip.ProxyHost, err)
			db.Rollback()
		}
	} else {
		UpdateIP(ipModel)
	}
	db.Commit()
}

// GetIPByProxyHost 根据 proxyHost 获取一条数据
func GetIPByProxyHost(host string) *IP {
	ipModel := new(IP)
	err := db.Model(&IP{}).Where("proxy_host = ?", host).Find(ipModel)
	if err.Error != nil {
		log.Errorf("get ip: %s, error msg: %v", ipModel.ProxyHost, err.Error)
		return nil
	}
	return ipModel
}

// CountIP 查询共有多少条数据
func CountIP() int64 {
	var count int64
	err := db.Model(&IP{}).Count(&count)
	if err.Error != nil {
		log.Errorf("ip count: %d, error msg: %v", count, err.Error)
		return -1
	}
	return count
}

// GetAllIP 获取所有数据
func GetAllIP() []*IP {
	var list []*IP
	err := GetDB().Model(&IP{}).Find(&list)
	ipCount := len(list)
	if err.Error != nil {
		log.Warnf("ip count: %d, error msg: %v\n", ipCount, err.Error)
		return nil
	}
	return list
}

// GetIPByProxyType 根据 proxyType 获取一条数据
func GetIPByProxyType(proxyType string) ([]IP, error) {
	list := make([]IP, 0)
	err := db.Model(&IP{}).Where("proxy_type = ?", proxyType).Find(&list)
	if err.Error != nil {
		log.Errorf("error msg: %v\n", err.Error)
		return list, err.Error
	}
	return list, nil
}

// UpdateIP 更新数据
func UpdateIP(ip *IP) {
	db := GetDB().Begin()
	ipModel := ip
	ipMap := make(map[string]interface{}, 0)
	ipMap["proxy_speed"] = ip.ProxySpeed
	ipMap["proxy_type"] = strings.ToLower(ip.ProxyType)

	ipMap["update_time"] = time.Now()
	if ipModel.ProxyId != 0 {
		err := db.Model(&IP{}).Where("proxy_id = ?", ipModel.ProxyId).Updates(ipMap)
		if err.Error != nil {
			log.Errorf("update ip: %s, error msg: %v", ipModel.ProxyHost, err.Error)
			db.Rollback()
		}
	}
	db.Commit()
}

// DeleteIP 删除数据
func DeleteIP(ip *IP) {
	db := GetDB().Begin()
	ipModel := ip
	err := db.Model(&IP{}).Where("proxy_id = ?", ipModel.ProxyId).Delete(ipModel)
	if err.Error != nil {
		log.Errorf("delete ip: %s, error msg: %v", ipModel.ProxyHost, err.Error)
		db.Rollback()
	}
	db.Commit()
}
func DeleteByIP(ip string) {
	db := GetDB().Begin()
	ipModel := IP{ProxyHost: ip}
	err := db.Model(&IP{}).Where("proxy_proxy_host = ?", ip).Delete(ipModel)
	if err.Error != nil {
		log.Errorf("delete ip: %s, error msg: %v", ipModel.ProxyHost, err.Error)
		db.Rollback()
	}
	db.Commit()
}

func Count() map[string]int64 {
	countMap := make(map[string]int64)
	proxyTypes := []string{"http", "https", "tcp", "socks5", "socks4", "tcp"}
	for _, proxyType := range proxyTypes {
		var count int64
		if err := GetDB().Model(&IP{}).Where("proxy_type = ?", proxyType).Count(&count).Error; err != nil {
			log.Errorf("count proxy_type: %s, error %s", proxyType, err)
		}
		countMap[proxyType] = count
	}
	return countMap
}
