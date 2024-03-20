package config

import (
	"time"

	"github.com/baowk/dilu-rd/config"
)

// type AppConfig struct {
// 	App AppInfo `mapstructure:"app" json:"app" yaml:"app"`
// }

type AppConfig struct {
	Server    ServerConf `mapstructure:"server" json:"server" yaml:"server"`
	RemoteCfg RemoteCfg  `mapstructure:"remote-cfg" json:"remote-cfg" yaml:"remote-cfg"`
	Rules     []RuleConf `mapstructure:"rules" json:"rules" yaml:"rules"`
	//Logger    LogCfg        `mapstructure:"logger" json:"logger" yaml:"logger"`
	JWT      JWTConf       `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	RdConfig config.Config `mapstructure:"rd-config" json:"rd-config" yaml:"rd-config"`
	Extend   Extend        `mapstructure:"extend" json:"extend" yaml:"extend"`
}

type ServerConf struct {
	Mode         string `mapstructure:"mode" json:"mode" yaml:"mode"`
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Name         string `mapstructure:"name" json:"name" yaml:"name"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	ReadTimeout  int    `mapstructure:"readtimeout" json:"readtimeout" yaml:"readtimeout"`
	WriteTimeout int    `mapstructure:"writertimeout" json:"writetimeout" yaml:"writetimeout"`
}

type RuleConf struct {
	Rd        bool     `mapstructure:"rd" json:"rd" yaml:"rd"`
	RuleName  string   `mapstructure:"name" json:"name" yaml:"name"`
	Prefix    string   `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Upstreams []string `mapstructure:"upstreams" json:"upstreams" yaml:"upstreams"`
	Rewrite   string   `mapstructure:"rewrite" json:"rewrite" yaml:"rewrite"`
	Handlers  []string `mapstructure:"handlers" json:"handlers" yaml:"handlers"`
}

type JWTConf struct {
	Secret  string `mapstructure:"secret"`
	Timeout int64  `mapstructure:"timeout"`
	Refresh int    `mapstructure:"refresh" json:"refresh" yaml:"refresh"` // 刷新时长
	Issuer  string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`    // 签发人
	Subject string `mapstructure:"subject" json:"subject" yaml:"subject"` // 签发主体
}

type RemoteCfg struct {
	Enable        bool          `mapstructure:"enable" json:"enable" yaml:"enable"`                         //是否开启远程配置
	Provider      string        `mapstructure:"provider" json:"provider" yaml:"provider"`                   //提供方
	Endpoint      string        `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`                   //端点
	Path          string        `mapstructure:"path" json:"path" yaml:"path"`                               //路径
	SecretKeyring string        `mapstructure:"secret-keyring" json:"secret-keyring" yaml:"secret-keyring"` //安全
	ConfigType    string        `mapstructure:"config-type" json:"config-type" yaml:"config-type"`          //配置类型
	Duration      time.Duration `mapstructure:"duration" json:"duration" yaml:"duration"`                   //重试时长
}

func (e *RemoteCfg) GetDuration() time.Duration {
	if e.Duration < 0 {
		return time.Second * 10
	}
	return e.Duration
}

func (e *RemoteCfg) GetConfigType() string {
	if e.ConfigType == "" {
		return "yaml"
	}
	return e.ConfigType
}

type Extend struct {
	Auth Auth
}

type Auth struct {
	BaseUrl string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
}
