package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Generate 生成JWT Token
func Generate(claims jwt.Claims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成签名字符串
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// AdminCustomClaims 自定义格式内容
type CustomClaims struct {
	UserId               int    `json:"uid,omitempty"`
	RoleId               int    `json:"rid,omitempty"`
	Phone                string `json:"mob,omitempty"`
	Nickname             string `json:"nick,omitempty"`
	JwtData              map[string]any
	jwt.RegisteredClaims // 内嵌标准的声明
}

func (c *CustomClaims) AddData(key string, val any) *CustomClaims {
	if c.JwtData == nil {
		c.JwtData = make(map[string]any, 0)
	}
	c.JwtData[key] = val
	return c
}

// func (c *CustomClaims) GetInt(key string) int {
// 	if val, ok := c.JwtData[key]; ok {
// 		return utils.GetInterfaceToInt(val)
// 	}
// 	return 0
// }

func (c *CustomClaims) GetString(key string) string {
	if val, ok := c.JwtData[key]; ok {
		return fmt.Sprintf("%s", val)
	}
	return ""
}

func (c *CustomClaims) ExpiresAt(expiresAt time.Time) *CustomClaims {
	c.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	return c
}

// NewAdminCustomClaims 初始化AdminCustomClaims
func NewClaims(userId int, expiresAt time.Time, issuer, subject string) CustomClaims {
	//now := time.Now()
	return CustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 定义过期时间
			Issuer:    issuer,                        // 签发人
			//IssuedAt:  jwt.NewNumericDate(now),       // 签发时间
			Subject: subject, // 签发主体
			//NotBefore: jwt.NewNumericDate(now),       // 生效时间
		},
	}
}
