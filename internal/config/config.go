package config

import (
	"Gin-K8S-Client/etc"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"strconv"
)

// keyMap yaml属性映射
var keyMap map[KeyName]string

type Config struct {
	Server Server
}

type Server struct {
	Host string `mapstructure:"host"`
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}

// init 初始化
func init() {
	var config *Config
	//viper相关配置
	viper.AddConfigPath(etc.YamlPath)
	viper.SetConfigName(etc.YamlName)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		klog.Fatal("config.init can't resolve config.yaml: " + err.Error())
		return
	}
	//映射
	if err := viper.Unmarshal(&config); err != nil {
		klog.Fatal("config.init can't unmarshal to config: " + err.Error())
		return
	}
	keyMap = make(map[KeyName]string)
	//初始化keyMap
	keyMap[ServerHost] = config.Server.Host
	keyMap[ServerPort] = config.Server.Port
	keyMap[ServerName] = config.Server.Name
}

func GetString(key KeyName) string {
	if v, ok := keyMap[key]; ok {
		return v
	}
	return ""
}

func GetInt(key KeyName) int {
	if v, ok := keyMap[key]; ok {
		res, err := strconv.Atoi(v)
		if err != nil {
			klog.Fatal("config.GetInt can't convert to int: " + err.Error())
			return -1
		}
		return res
	}
	return -1
}
