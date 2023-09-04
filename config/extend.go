package config

// Extend 扩展配置
//
//	extend:
//	  demo:
//	    name: demo-name

type Extend struct {
	Auth Auth
}

type Auth struct {
	BaseUrl string
}
