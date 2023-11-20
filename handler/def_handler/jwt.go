package def_handler

import (
	"dilu-gateway/handler/def_handler/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	Secret = "S23456789" //密码自行设定
	// TokenExpireDuration = appConfig.JwtConfig.Timeout * int64(time.Second) //设置过期时间

	TokenLookup   = "header: Authorization, query: token, cookie: token"
	TokenHeadName = "Bearer"

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cokie is empty
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrEmptyParamToken can be thrown if authing with parameter in path, the parameter in path is empty
	ErrEmptyParamToken = errors.New("parameter token is empty")
)

type JwtProxyHandler struct {
	Secret    string
	ExpiresAt int64

	//	Pattern    *regexp.Regexp//"^/v2/.+/auth/.*$"
	HeaderKey  string
	HeaderName string
	QueryKey   string
	CookieKey  string
	Refresh    int
	Issuer     string
	Subject    string
}

type JwtProxyHandlerBuilder struct {
	h JwtProxyHandler
}

func NewJwt() JwtProxyHandlerBuilder {
	return JwtProxyHandlerBuilder{}
}

func (b JwtProxyHandlerBuilder) Refresh(refresh int) JwtProxyHandlerBuilder {
	b.h.Refresh = refresh
	return b
}

func (b JwtProxyHandlerBuilder) Subject(subject string) JwtProxyHandlerBuilder {
	b.h.Subject = subject
	return b
}

func (b JwtProxyHandlerBuilder) Issuer(issuer string) JwtProxyHandlerBuilder {
	b.h.Issuer = issuer
	return b
}

func (b JwtProxyHandlerBuilder) Secret(secret string) JwtProxyHandlerBuilder {
	b.h.Secret = secret
	return b
}

func (b JwtProxyHandlerBuilder) ExpiresAt(expiresAt int64) JwtProxyHandlerBuilder {
	b.h.ExpiresAt = expiresAt
	return b
}

func (b JwtProxyHandlerBuilder) HeaderKey(headerKey string) JwtProxyHandlerBuilder {
	b.h.HeaderKey = headerKey
	return b
}

func (b JwtProxyHandlerBuilder) HeaderName(headerName string) JwtProxyHandlerBuilder {
	b.h.HeaderName = headerName
	return b
}

func (b JwtProxyHandlerBuilder) QueryKey(queryKey string) JwtProxyHandlerBuilder {
	b.h.QueryKey = queryKey
	return b
}

func (b JwtProxyHandlerBuilder) CookieKey(cookieKey string) JwtProxyHandlerBuilder {
	b.h.CookieKey = cookieKey
	return b
}

func (b JwtProxyHandlerBuilder) Build() JwtProxyHandler {
	if b.h.HeaderKey == "" && b.h.QueryKey == "" && b.h.CookieKey == "" {
		b.h.HeaderKey = "Authorization"
		b.h.HeaderName = "Bearer"
	}
	return b.h
}

func (h JwtProxyHandler) BeforeHander(w http.ResponseWriter, r *http.Request, args ...interface{}) (int, string) {
	var tokenStr string
	var err error
	if h.HeaderKey != "" {
		tokenStr, err = jwtFromHeader(r, h.HeaderKey, h.HeaderName)
		if err != nil && err == ErrInvalidAuthHeader {
			return 1001, "Token有误"
		}
	}

	if len(tokenStr) == 0 && h.QueryKey != "" {
		tokenStr, _ = jwtFromQuery(r, h.QueryKey)
	}

	if len(tokenStr) == 0 && h.CookieKey != "" {
		tokenStr, _ = jwtFromCookie(r, "token")
	}
	if len(tokenStr) == 0 {
		r.Header.Del("userId")
		return 200, "未找到Token"
	}
	customClaims := new(utils.CustomClaims)
	err = ParseToken(tokenStr, customClaims, h.Secret, jwt.WithSubject(h.Subject))
	if err != nil || customClaims == nil {
		return 401, "Token有误"
	}

	exp, err := customClaims.GetExpirationTime()
	// 获取过期时间返回err,或者exp为nil都返回错误
	if err != nil || exp == nil {
		return 401, "token已失效"
	}

	// 刷新时间大于0则判断剩余时间小于刷新时间时刷新Token并在Response header中返回
	if h.Refresh > 0 {
		now := time.Now()
		diff := exp.Time.Sub(now)
		refreshTTL := time.Duration(h.Refresh) * time.Minute
		//fmt.Println(diff.Seconds(), refreshTTL)
		if diff < refreshTTL {
			exp := time.Now().Add(time.Duration(h.ExpiresAt) * time.Minute)
			customClaims.ExpiresAt(exp)
			newToken, _ := Refresh(customClaims, h.Secret)
			r.Header.Set("refresh-access-token", newToken)
			r.Header.Set("refresh-exp", strconv.FormatInt(exp.Unix(), 10))
		}
	}

	r.Header.Set("a_uid", fmt.Sprintf("%d", customClaims.UserId))
	r.Header.Set("a_rid", fmt.Sprintf("%d", customClaims.RoleId))
	r.Header.Set("a_mobile", customClaims.Phone)
	r.Header.Set("a_nickname", customClaims.Nickname)
	//r.Header.Set("jwt_data", customClaims.JwtData)
	return 200, ""
}

func Refresh(claims jwt.Claims, secretKey string) (string, error) {
	return utils.Generate(claims, secretKey)
}

func (h JwtProxyHandler) AfferHandler(w http.ResponseWriter, r *http.Request, args ...interface{}) (int, string) {
	fmt.Println("AFT JwtProxyHandler")
	return 200, ""
}

func (h JwtProxyHandler) GetName() string {
	return "jwt"
}

func ParseToken(tokenString string, claims jwt.Claims, secret string, options ...jwt.ParserOption) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secret), err
	}, options...)
	if err != nil {
		return err
	}

	// 对token对象中的Claim进行类型断言
	if token.Valid { // 校验token
		return nil
	}
	return errors.New("Invalid Token")
}

func jwtFromHeader(r *http.Request, key, name string) (string, error) {
	authHeader := r.Header.Get(key)
	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == name) {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func jwtFromQuery(r *http.Request, key string) (string, error) {
	token := r.FormValue(key)
	if token == "" {
		return "", ErrEmptyQueryToken
	}

	return token, nil
}

func jwtFromCookie(r *http.Request, key string) (string, error) {
	if cookie, err := r.Cookie(key); err == nil {
		return cookie.Value, nil
	}
	return "", nil
}

// func GetInterfaceToInt(t1 interface{}) int {
// 	var t2 int
// 	switch t1.(type) {
// 	case uint:
// 		t2 = int(t1.(uint))
// 		break
// 	case int8:
// 		t2 = int(t1.(int8))
// 		break
// 	case uint8:
// 		t2 = int(t1.(uint8))
// 		break
// 	case int16:
// 		t2 = int(t1.(int16))
// 		break
// 	case uint16:
// 		t2 = int(t1.(uint16))
// 		break
// 	case int32:
// 		t2 = int(t1.(int32))
// 		break
// 	case uint32:
// 		t2 = int(t1.(uint32))
// 		break
// 	case int64:
// 		t2 = int(t1.(int64))
// 		break
// 	case uint64:
// 		t2 = int(t1.(uint64))
// 		break
// 	case float32:
// 		t2 = int(t1.(float32))
// 		break
// 	case float64:
// 		t2 = int(t1.(float64))
// 		break
// 	case string:
// 		t2, _ = strconv.Atoi(t1.(string))
// 		if t2 == 0 && len(t1.(string)) > 0 {
// 			f, _ := strconv.ParseFloat(t1.(string), 64)
// 			t2 = int(f)
// 		}
// 		break
// 	case nil:
// 		t2 = 0
// 		break
// 	case json.Number:
// 		t3, _ := t1.(json.Number).Int64()
// 		t2 = int(t3)
// 		break
// 	default:
// 		t2 = t1.(int)
// 		break
// 	}
// 	return t2
// }
