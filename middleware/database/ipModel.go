package database

import (
	"IpProxyPool/util"
	"github.com/sirupsen/logrus"
	"strings"
)

// IP struct
type IP struct {
	ProxyId       int64  `gorm:"primary_key; auto_increment; not null" json:"-"`
	ProxyHost     string `gorm:"type:varchar(255); not null; uniqueIndex:UNIQUE_HOST_PORT" json:"proxyHost"`
	ProxyPort     int    `gorm:"type:int(11); not null; uniqueIndex:UNIQUE_HOST_PORT" json:"proxyPort"`
	ProxyType     string `gorm:"type:varchar(64); not null" json:"proxyType"`
	ProxyLocation string `gorm:"type:varchar(255); default null" json:"proxyLocation"`
	ProxySpeed    int    `gorm:"type:int(20); not null; default 0" json:"proxySpeed"`
	ProxySource   string `gorm:"type:varchar(64); not null;" json:"proxySource"`
	CreateTime    string `gorm:"type:varchar(50); not null" json:"-"`
	UpdateTime    string `gorm:"type:varchar(50); default ''" json:"updateTime"`
}

func (i *IP) TableName() string {
	return "ip"
}

// SaveIp 保存数据到数据库
func SaveIp(ip *IP) {
	db := GetDB().Begin()
	ip.ProxyType = strings.ToLower(ip.ProxyType)
	ipModel := GetIpByProxyHost(ip.ProxyHost)
	if ipModel.ProxyHost == "" {
		err := db.Model(new(IP)).Create(ip)
		if err.Error != nil {
			logrus.Errorf("save ip: %s, error msg: %v", ip.ProxyHost, err.Error)
			db.Rollback()
		}
	} else {
		UpdateIp(ipModel)
	}
	db.Commit()
}

// GetIpByProxyHost 根据 proxyHost 获取一条数据
func GetIpByProxyHost(host string) *IP {
	db := GetDB()
	ipModel := new(IP)
	err := db.Model(new(IP)).Where("proxy_host = ?", host).Find(ipModel)
	if err.Error != nil {
		logrus.Errorf("get ip: %s, error msg: %v", ipModel.ProxyHost, err.Error)
		return nil
	}
	return ipModel
}

// CountIp 查询共有多少条数据
func CountIp() int64 {
	db := GetDB()
	var count int64
	err := db.Model(new(IP)).Count(&count)
	if err.Error != nil {
		logrus.Errorf("ip count: %d, error msg: %v", count, err.Error)
		return -1
	}
	return count
}

// GetAllIp 获取所有数据
func GetAllIp() []IP {
	db := GetDB()
	list := make([]IP, 0)
	err := db.Model(new(IP)).Find(&list)
	ipCount := len(list)
	if err.Error != nil {
		logrus.Warnf("ip count: %d, error msg: %v\n", ipCount, err.Error)
		return nil
	}
	return list
}

// GetIpByProxyType 根据 proxyType 获取一条数据
func GetIpByProxyType(proxyType string) ([]IP, error) {
	db := GetDB()
	list := make([]IP, 0)
	err := db.Model(new(IP)).Where("proxy_type = ?", proxyType).Find(&list)
	if err.Error != nil {
		logrus.Errorf("error msg: %v\n", err.Error)
		return list, err.Error
	}
	return list, nil
}

// UpdateIp 更新数据
func UpdateIp(ip *IP) {
	db := GetDB().Begin()
	ipModel := ip
	ipMap := make(map[string]interface{}, 0)
	ipMap["proxy_speed"] = ip.ProxySpeed
	ipMap["proxy_type"] = ip.ProxyType

	ipMap["update_time"] = util.FormatDateTime()
	if ipModel.ProxyId != 0 {
		err := db.Model(new(IP)).Where("proxy_id = ?", ipModel.ProxyId).Updates(ipMap)
		if err.Error != nil {
			logrus.Errorf("update ip: %s, error msg: %v", ipModel.ProxyHost, err.Error)
			db.Rollback()
		}
	}
	db.Commit()
}

// DeleteIp 删除数据
func DeleteIp(ip *IP) {
	db := GetDB().Begin()
	ipModel := ip
	err := db.Model(new(IP)).Where("proxy_id = ?", ipModel.ProxyId).Delete(ipModel)
	if err.Error != nil {
		logrus.Errorf("delete ip: %s, error msg: %v", ipModel.ProxyHost, err.Error)
		db.Rollback()
	}
	db.Commit()
}

type ProxyTypeCount struct {
	Http  int64 `json:"http"`
	Https int64 `json:"https"`
	Other int64 `json:"other"`
}

func Count() (c *ProxyTypeCount) {
	c = &ProxyTypeCount{}

	err := GetDB().Model(new(IP)).Where("proxy_type = ?", "http").Count(&c.Http).Error
	if err != nil {
		logrus.Errorf("error msg: %v\n", err.Error())
	}

	err = GetDB().Model(new(IP)).Where("proxy_type = ?", "https").Count(&c.Https).Error
	if err != nil {
		logrus.Errorf("error msg: %v\n", err.Error())
	}
	err = GetDB().Model(new(IP)).Where("proxy_type != ? and proxy_type != ? ", "http", "https").Count(&c.Other).Error
	if err != nil {
		logrus.Errorf("error msg: %v\n", err.Error())
	}

	return
}
