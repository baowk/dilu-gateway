package config

// type AppConfig struct {
// 	App AppInfo `mapstructure:"app" json:"app" yaml:"app"`
// }

type AppInfo struct {
	Server ServerConf        `mapstructure:"server" json:"server" yaml:"server"`
	Rules  []RuleConf        `mapstructure:"rules" json:"rules" yaml:"rules"`
	Logger LogCfg            `mapstructure:"logger" json:"logger" yaml:"logger"`
	JWT    JWTConf           `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Extend map[string]string `mapstructure:"extend" json:"extend" yaml:"extend"`
}

type ServerConf struct {
	Mode         string `mapstructure:"mode" json:"mode" yaml:"mode"`
	RemoteConfig bool   `mapstructure:"remoteconfig" json:"remoteconfig" yaml:"remoteconfig"`
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Name         string `mapstructure:"name" json:"name" yaml:"name"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	ReadTimeout  int    `mapstructure:"readtimeout" json:"readtimeout" yaml:"readtimeout"`
	WriteTimeout int    `mapstructure:"writertimeout" json:"writetimeout" yaml:"writetimeout"`
}

type RuleConf struct {
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

// type DatabaseConf struct {
// 	Driver string `mapstructure:"driver"`
// 	Source string `mapstructure:"source"`
// }

// type CacheConf struct {
// 	Redis RedisConf `mapstructure:"redis"`
// }

// type RedisConf struct {
// 	Addr     string `mapstructure:"addr"`
// 	Password string `mapstructure:"password"`
// 	DB       int    `mapstructure:"db"`
// }

// func LoadFile(filePath string) AppInfo {
// 	viper.SetConfigFile(filePath) // 指定配置文件路径
// 	//viper.SetEnvPrefix("app")
// 	// 读取配置文件
// 	if err := viper.ReadInConfig(); err != nil {
// 		panic(err) // 读取配置文件失败
// 	}
// 	var config AppInfo

// 	if err := viper.Unmarshal(&config); err != nil {
// 		log.Fatal(err) // 解析配置文件失败
// 	}

// 	//fmt.Println(viper.AllKeys())
// 	return config
// }
