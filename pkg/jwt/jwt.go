package jwt

import (
	"errors"
	"mianshiba/conf"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

const TokenKey = "jwt_key"

// JWTClaims 自定义JWT声明结构
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID int64, username, role string) (string, error) {
	cfg := conf.Global.Security

	// 创建声明
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWTExpirationInt) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "interview-agent",
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串token
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// parseToken 解析JWT token
func ParseToken(jwtToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(conf.Global.Security.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserID 从上下文获取用户ID
func GetUserID(ctx *app.RequestContext) uint {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetUsername 从上下文获取用户名
func GetUsername(ctx *app.RequestContext) string {
	username, exists := ctx.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

// GetUserRole 从上下文获取用户角色
func GetUserRole(ctx *app.RequestContext) string {
	role, exists := ctx.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}
