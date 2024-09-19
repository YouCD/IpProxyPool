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
	ProxyId       int64     `gorm:"column:proxy_id;primaryKey;autoIncrement;not null" json:"proxyId"`
	ProxyHost     string    `gorm:"column:proxy_host;type:varchar(255);not null;uniqueIndex:UNIQUE_HOST_PORT,uniqueIndexLength:191" json:"proxyHost"`
	ProxyPort     int       `gorm:"column:proxy_port;type:int;not null;uniqueIndex:UNIQUE_HOST_PORT" json:"proxyPort"`
	ProxyType     string    `gorm:"column:proxy_type;type:varchar(64);not null" json:"proxyType"`
	ProxyLocation string    `gorm:"column:proxy_location;type:varchar(255);not null" json:"proxyLocation"`
	ProxySpeed    int       `gorm:"column:proxy_speed;type:int;not null;default:0" json:"proxySpeed"`
	ProxySource   string    `gorm:"column:proxy_source;type:varchar(64);not null" json:"proxySource"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;not null" json:"createTime"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;not null" json:"updateTime"`
}

func (i *IP) TableName() string {
	return "ip"
}

// SaveIP 保存数据到数据库
func SaveIP(ip *IP) {
	if ip.ProxyHost == "" {
		return
	}
	ip.ProxyType = strings.ToLower(ip.ProxyType)
	ipModel := GetIPByProxyHost(ip.ProxyHost)
	if ipModel.ProxyId != 0 {
		UpdateIP(ip)
		return
	}
	if err := db.Model(&IP{}).Create(ip).Error; err != nil {
		log.Errorf("save ip: %s, error msg: %v", ip.ProxyHost, err)
	}
}

// GetIPByProxyHost 根据 proxyHost 获取一条数据
func GetIPByProxyHost(host string) *IP {
	ipModel := &IP{}
	if err := db.Model(&IP{}).Where("proxy_host = ?", host).Scan(ipModel).Error; err != nil {
		log.Errorf("get ip: %s, error msg: %s", host, err)
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
	ipMap := make(map[string]interface{}, 0)
	ipMap["proxy_speed"] = ip.ProxySpeed
	ipMap["proxy_type"] = strings.ToLower(ip.ProxyType)
	ipMap["update_time"] = time.Now()
	if ip.ProxyId != 0 {
		tx := GetDB().Begin()
		if err := tx.Model(ip).Where("proxy_id = ?", ip.ProxyId).Updates(ipMap).Error; err != nil {
			log.Errorf("update ip: %s, error msg: %v", ip.ProxyHost, err)
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
}

// DeleteIP 删除数据
func DeleteIP(ip *IP) {
	tx := GetDB().Begin()
	ipModel := ip
	err := tx.Model(&IP{}).Where("proxy_id = ?", ipModel.ProxyId).Delete(ipModel)
	if err.Error != nil {
		log.Errorf("delete ip: %s, error msg: %v", ipModel.ProxyHost, err.Error)
		tx.Rollback()
	}
	tx.Commit()
}
func DeleteByIP(ip string) {
	tx := GetDB().Begin()
	ipModel := IP{}
	if err := tx.Model(&IP{}).Where("proxy_host = ?", ip).Scan(&ipModel).Delete(&ipModel).Error; err != nil {
		log.Errorf("delete ip: %s, error msg: %v", ipModel.ProxyHost, err)
		tx.Rollback()
	}
	tx.Commit()
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
