package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"dilu-gateway/config"
	"dilu-gateway/handler"
	"dilu-gateway/handler/def_handler"
)

type Upstream struct {
	Server   string
	Limit    int
	ok       bool
	downTime int
}

type Rule struct {
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

func Run(conf config.AppConfig) {
	Append(def_handler.NewJwt().Secret(conf.App.JWT.Secret).ExpiresAt(conf.App.JWT.Timeout).
		Subject(conf.App.JWT.Subject).Issuer(conf.App.JWT.Issuer).Refresh(conf.App.JWT.Refresh).Build())
	m := conf.App.Extend
	Append(def_handler.AuthProxyHandler{BaseURL: m["authbaseurl"]})
	for _, ruleC := range conf.App.Rules {
		rule := Rule{
			Name:      ruleC.RuleName,
			Prefix:    ruleC.Prefix,
			Upstreams: ruleC.Upstreams,
			Rewrite:   ruleC.Rewrite,
		}
		for _, hname := range ruleC.Handlers {
			if h, ok := handlerMap[hname]; ok {
				rule.Handlers = append(rule.Handlers, &h)
			}
		}
		rules = append(rules, rule)
	}
	// 创建一个自定义的请求处理程序
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// 根据规则匹配请求路径
		for _, rule := range rules {
			fmt.Printf("uri:%s,prefix:%s,匹配 %v\n", r.RequestURI, rule.Prefix, strings.HasPrefix(r.RequestURI, rule.Prefix))
			if strings.HasPrefix(r.RequestURI, rule.Prefix) {
				var upstream string
				size := len(rule.Upstreams)
				if size == 1 {
					upstream = rule.Upstreams[0]
				} else {
					upstream = rule.Upstreams[rand.Intn(size)]
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
					log.Println(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
				for _, handler := range rule.Handlers {
					code, msg := (*handler).BeforeHander(w, r)
					if code != 200 {
						data := map[string]interface{}{
							"code": code,
							"msg":  msg,
						}
						jsonBytes, err := json.Marshal(data)
						if err != nil {
							log.Fatal(err)
						}
						w.Write(jsonBytes)
						return
					}
				}

				targetURL.Path = ""

				fmt.Println("targetURL:" + targetURL.String() + "  path:" + targetURL.Path)
				// 创建代理，并将请求重定向到代理目标
				proxy := httputil.NewSingleHostReverseProxy(targetURL)

				proxy.ServeHTTP(w, r)
				fmt.Printf("用时：%v  \n", time.Now().Sub(startTime))
				return
			}
		}

		// 如果没有匹配的规则，则返回404 Not Found
		http.NotFound(w, r)

	})

	addr := fmt.Sprintf("%s:%d", conf.App.Server.Host, conf.App.Server.Port)
	fmt.Printf("启动服务监听%s \n", addr)
	// 启动HTTP服务器
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatal(err)
	}
}
