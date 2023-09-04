package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App AppInfo `mapstructure:"app"`
}

type AppInfo struct {
	Server   ServerConf        `mapstructure:"server"`
	Rules    []RuleConf        `mapstructure:"rules"`
	Logger   LoggerConf        `mapstructure:"logger"`
	JWT      JWTConf           `mapstructure:"jwt"`
	Database DatabaseConf      `mapstructure:"database"`
	Cache    CacheConf         `mapstructure:"cache"`
	Extend   map[string]string `mapstructure:"extend"`
}

type ServerConf struct {
	Mode         string `mapstructure:"mode"`
	Host         string `mapstructure:"host"`
	Name         string `mapstructure:"name"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"readtimeout"`
	WriteTimeout int    `mapstructure:"writertimeout"`
	EnabledP     bool   `mapstructure:"enabledp"`
}

type RuleConf struct {
	RuleName  string   `mapstructure:"name"`
	Prefix    string   `mapstructure:"prefix"`
	Upstreams []string `mapstructure:"upstreams"`
	Rewrite   string   `mapstructure:"rewrite"`
	Handlers  []string `mapstructure:"handlers"`
}

type LoggerConf struct {
	Path      string `mapstructure:"path"`
	Stdout    string `mapstructure:"stdout"`
	Level     string `mapstructure:"level"`
	EnableDDB bool   `mapstructure:"enableddb"`
}

type JWTConf struct {
	Secret  string `mapstructure:"secret"`
	Timeout int64  `mapstructure:"timeout"`
}

type DatabaseConf struct {
	Driver string `mapstructure:"driver"`
	Source string `mapstructure:"source"`
}

type CacheConf struct {
	Redis RedisConf `mapstructure:"redis"`
}

type RedisConf struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func LoadFile(filePath string) AppConfig {
	viper.SetConfigFile(filePath) // 指定配置文件路径
	//viper.SetEnvPrefix("app")
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(err) // 读取配置文件失败
	}
	var config AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err) // 解析配置文件失败
	}

	//fmt.Println(viper.AllKeys())
	return config
}
