package def_handler

import (
	"dilu-gateway/common"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AuthProxyHandler struct {
	BaseURL string
}

var client common.HTTPClient
var perms = make(map[string]map[string]int, 0)
var initf = false

func (h AuthProxyHandler) BeforeHander(w http.ResponseWriter, r *http.Request, args ...interface{}) (int, string) {
	if !initf {
		client = common.HTTPClient{
			BaseURL: h.BaseURL,
		}
		if err := getAllPerms(); err != nil {
			return 1025, "系统暂时无法服务"
		}
	}
	uri := r.Method + ":" + r.RequestURI
	if m, ok := perms[uri]; ok {
		uid := r.Header.Get("userId")
		if uid == "" {
			return 1023, "未登录，无法访问"
		}
		companyId := r.Header.Get("companyId")
		if companyId == "" {
			return 1023, "需选择企业后操作"
		}
		up := getUserPerms(companyId, uid)
		if up == "" {
			return 1024, "没有访问权限"
		}
		ids := strings.Split(up, ",")
		for _, id := range ids {
			if id != "" {
				if _, ok := m[id]; ok {
					return 200, ""
				}
			}
		}
		return 1023, "没有访问权限"
	}

	return 200, ""
}

func (h AuthProxyHandler) AfferHandler(w http.ResponseWriter, r *http.Request, args ...interface{}) (int, string) {
	fmt.Println("AFT AuthProxyHandler")
	return 200, ""
}

func (h AuthProxyHandler) GetName() string {
	return "auth"
}

var nilReqData = []byte(`{}`)

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func getAllPerms() error {
	resB, err := client.Post("/v2/team/perms", nilReqData)
	if err != nil {
		fmt.Print(err)
		return err
	}
	var res Res
	if err = json.Unmarshal(resB, &res); err != nil {
		return err
	}
	if res.Code == 200 {
		d, err := json.Marshal(res.Data)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(d, &perms); err == nil {
			initf = true
		} else {
			return err
		}
	} else {
		return errors.New(res.Msg)
	}
	return nil
}

func getUserPerms(companyId, uid string) string {
	client.AddHeader("userId", uid)
	client.AddHeader("companyId", companyId)
	resB, err := client.Post("/v2/team/myPerms", nilReqData)
	if err != nil {
		return ""
	}
	var res Res
	if err = json.Unmarshal(resB, &res); err != nil {
		return ""
	}
	if res.Code == 200 {
		d, err := json.Marshal(res.Data)
		if err != nil {
			return ""
		}
		return string(d)
	}
	return ""
}
