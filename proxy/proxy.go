package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"dilu-gateway/config"
	"dilu-gateway/handler"
	"dilu-gateway/handler/def"

	"github.com/baowk/dilu-rd/rd"
)

var (
	Cfg *config.AppConfig
	rdc rd.RDClient
)

type Upstream struct {
	Server   string
	Limit    int
	ok       bool
	downTime int
}

type Rule struct {
	Rd        bool
	Upstream  string
	Name      string                  //名称
	Prefix    string                  //匹配前缀
	Upstreams []string                //后端 服务器
	Rewrite   string                  //重写
	Handlers  []*handler.ProxyHandler //处理器
	//pattern   *regexp.Regexp
}

var rules = make([]Rule, 0)

var handlerMap = make(map[string]handler.ProxyHandler, 0)

func Append(h handler.ProxyHandler) {
	handlerMap[h.GetName()] = h
}

func Run() {
	InitRd()
	InitHandler()
	// 创建一个自定义的请求处理程序
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// 根据规则匹配请求路径
		for _, rule := range rules {
			if strings.HasPrefix(r.RequestURI, rule.Prefix) {
				var upstream string
				if rule.Rd {
					// 从注册中心获取服务列表
					node, err := rdc.GetService(rule.Name, r.RemoteAddr)
					if err != nil {
						slog.Error("get service error", err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
					upstream = fmt.Sprintf("%s://%s:%d", node.Protocol, node.Addr, node.Port)
				} else {
					size := len(rule.Upstreams)
					if size == 1 {
						upstream = rule.Upstreams[0]
					} else {
						upstream = rule.Upstreams[rand.Intn(size)]
					}
				}

				var tgUrl string
				if len(rule.Rewrite) == 0 {
					tgUrl = upstream + r.RequestURI
				} else {
					tgUrl = upstream + strings.Replace(r.RequestURI, rule.Prefix, rule.Rewrite, 1)
				}
				// 解析代理目标URL
				targetURL, err := url.Parse(tgUrl)
				if err != nil {
					slog.Error("parse error", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				for _, handler := range rule.Handlers {
					code, msg := (*handler).BeforeHander(w, r)
					if code != 200 {
						data := map[string]interface{}{
							"code": code,
							"msg":  msg,
						}
						jsonBytes, err := json.Marshal(data)
						if err != nil {
							slog.Error("marshal err", err)
						}
						w.Write(jsonBytes)
						slog.Warn("before", "url", tgUrl, "handler", (*handler).GetName(), "msg", msg)
						return
					}
				}

				targetURL.Path = ""
				r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
				//slog.Debug("target", "URL", targetURL.String()), "path", targetURL.Path))
				// 创建代理，并将请求重定向到代理目标
				proxy := httputil.NewSingleHostReverseProxy(targetURL)
				proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
					slog.Error("proxy error", "url", tgUrl, err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				proxy.ServeHTTP(w, r)
				slog.Debug("times", tgUrl, time.Now().Sub(startTime))
				return
			}
		}

		// 如果没有匹配的规则，则返回404 Not Found
		http.NotFound(w, r)
		slog.Error("no match", "RequestURI", r.RequestURI)

	})

	addr := fmt.Sprintf("%s:%d", Cfg.Server.Host, Cfg.Server.Port)
	fmt.Printf("启动服务监听:%s\n", addr)
	// 启动HTTP服务器
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatal(err)
	}
}

// func Initslog() {
// 	// 初始化日志
// 	slogInit()
// 	// 初始化注册中心

// }

func InitRd() {
	if Cfg.RdConfig.Enable {
		var err error
		rdc, err = rd.NewRDClient(&Cfg.RdConfig)
		if err != nil {
			log.Fatal("init rdclient error", err)
		}
	}
}

func InitHandler() {
	jwt := &def.JwtProxyHandler{
		ExpiresAt: Cfg.JWT.Timeout,
		Refresh:   Cfg.JWT.Refresh,
		Issuer:    Cfg.JWT.Issuer,
		Subject:   Cfg.JWT.Subject,
		Secret:    Cfg.JWT.Secret,
	}
	jwt.Build()
	Append(jwt)
	Append(&def.AuthProxyHandler{BaseURL: Cfg.Extend.Auth.BaseUrl})

	for _, ruleC := range Cfg.Rules {
		rule := Rule{
			Name:      ruleC.RuleName,
			Prefix:    ruleC.Prefix,
			Upstreams: ruleC.Upstreams,
			Rewrite:   ruleC.Rewrite,
			Rd:        ruleC.Rd,
		}
		for _, hname := range ruleC.Handlers {
			if h, ok := handlerMap[hname]; ok {
				rule.Handlers = append(rule.Handlers, &h)
			}
		}
		rules = append(rules, rule)
	}
}
