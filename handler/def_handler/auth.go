package def_handler

import (
	"dilu-gateway/common"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type AuthProxyHandler struct {
	BaseURL string
}

var client common.HTTPClient
var perms = make(map[string]int, 0)
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
		uid := r.Header.Get("a_uid")
		if uid == "" {
			return 1023, "未登录，无法访问"
		}
		teamId := r.Header.Get("team_id")
		if teamId == "" {
			return 1023, "需选择团队后操作"
		}
		if getUserPerms(m) != nil {
			return 1023, "没有访问权限"
		}
	}

	return 200, ""
}

func (h AuthProxyHandler) AfferHandler(w http.ResponseWriter, r *http.Request, args ...interface{}) (int, string) {
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

type SysApi struct {
	Id       int    `json:"id"`       //主键
	Title    string `json:"title"`    //标题
	Method   string `json:"method"`   //请求类型
	Path     string `json:"path"`     //请求地址
	PermType int    `json:"permType"` //权限类型（1：无需认证 2:须token 3：须鉴权）
	Status   int    `json:"status"`   //状态 3 DEF 2 OK 1 del
}

func getAllPerms() error {
	resB, err := client.Post("/api/v1/sys/sys-api/all", nilReqData)
	if err != nil {
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
		var list []SysApi
		if err := json.Unmarshal(d, &list); err == nil {
			for _, item := range list {
				perms[item.Method+":"+item.Path] = item.Id
			}
			initf = true
		} else {
			return err
		}
	} else {
		return errors.New(res.Msg)
	}
	return nil
}

func getUserPerms(m int) error {
	client.AddHeader("apid_id", strconv.Itoa(m))
	resB, err := client.Post("/api/v1/sys/canAccess", nilReqData)
	if err != nil {
		return err
	}
	var res Res
	if err = json.Unmarshal(resB, &res); err != nil {
		return err
	}
	if res.Code == 200 {
		return nil
	}
	return errors.New(res.Msg)
}
