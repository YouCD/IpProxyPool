package config

import (
	"IpProxyPool/util/fileutil"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/youcd/toolkit/log"
	"os"
	"path"
	"strings"
)

type System struct {
	AppName  string `yaml:"appName"`
	HttpAddr string `yaml:"httpAddr"`
	HttpPort string `yaml:"httpPort"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DbName   string `yaml:"dbName"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Log struct {
	FilePath string `yaml:"filePath"`
	FileName string `yaml:"fileName"`
	Level    string `yaml:"level"`
	Mode     string `yaml:"mode"`
}

type YamlSetting struct {
	System   System   `yaml:"system"`
	Database Database `yaml:"database"`
	Log      Log      `yaml:"log"`
}

var (
	Vip           = viper.New()
	ConfigFile    = ""
	ServerSetting = new(YamlSetting)
)

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	if ConfigFile != "" {
		if !fileutil.PathExists(ConfigFile) {
			log.Errorf("No such file or directory: %s", ConfigFile)
			os.Exit(-1)
		} else {
			// Use config file from the flag.
			Vip.SetConfigFile(ConfigFile)
			Vip.SetConfigType("yaml")
		}
	} else {
		log.Errorf("Could not find config file: %s", ConfigFile)
		os.Exit(-1)
	}
	// If a config file is found, read it in.
	err := Vip.ReadInConfig()
	if err != nil {
		log.Errorf("Failed to get config file: %s", ConfigFile)
	}
	Vip.WatchConfig()
	Vip.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s\n", e.Name)
		fmt.Printf("Config file changed: %s\n", e.Name)
		ServerSetting = GetConfig(Vip)
		log.SetLogLevel(ServerSetting.Log.Level)
		switch strings.ToLower(ServerSetting.Log.Mode) {
		case "file":
			log.Init(false)
			log.Info(path.Join(ServerSetting.Log.FilePath, ServerSetting.Log.FileName))
			log.SetFileName(path.Join(ServerSetting.Log.FilePath, ServerSetting.Log.FileName))
			log.SetLogLevel(ServerSetting.Log.Level)
		default:
			log.Init(true)
			log.SetLogLevel(ServerSetting.Log.Level)
		}
	})
	Vip.AllSettings()
	ServerSetting = GetConfig(Vip)
	log.Init(true)
	log.SetLogLevel(ServerSetting.Log.Level)
	log.SetFileName(path.Join(ServerSetting.Log.FilePath, ServerSetting.Log.FileName))
}

// 解析配置文件，反序列化
func GetConfig(vip *viper.Viper) *YamlSetting {
	setting := new(YamlSetting)
	// 解析配置文件，反序列化
	if err := vip.Unmarshal(setting); err != nil {
		log.Errorf("Unmarshal yaml faild: %s", err)
		os.Exit(-1)
	}
	return setting
}
