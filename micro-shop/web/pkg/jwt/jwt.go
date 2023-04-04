package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/wx-up/coding/micro-shop/web/config"
)

var (
	ErrTokenValid   = errors.New("令牌错误")
	ErrTokenExpired = errors.New("令牌已经过期")
)

type CustomClaims struct {
	Id       uint64
	Nickname string
	jwt.StandardClaims
}

type Manager struct {
	secret     string
	expireTime int64
}

type Option func(*Manager)

func WithSecret(s string) Option {
	return func(manager *Manager) {
		manager.secret = s
	}
}

func WithExpireTime(t int64) Option {
	return func(manager *Manager) {
		manager.expireTime = t
	}
}

func New(opts ...Option) *Manager {
	m := &Manager{
		secret:     config.Config().Jwt.Secret,
		expireTime: config.Config().Jwt.ExpireTime,
	}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

type CustomInfo struct {
	Id       uint64
	Nickname string
}

// Generate 生成 token
func (m *Manager) Generate(info CustomInfo) (string, error) {
	claims := &CustomClaims{
		Id:       info.Id,
		Nickname: info.Nickname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(m.expireTime) * time.Second).Unix(),
			NotBefore: time.Now().Unix(), // 生效时间
		},
	}
	res := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return res.SignedString(m.secret)
}

// ParseToken 解析 token
func (m *Manager) ParseToken(token string) (CustomInfo, error) {
	res, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		if validErr, ok := err.(*jwt.ValidationError); ok {
			if validErr.Errors&jwt.ValidationErrorExpired > 0 {
				return CustomInfo{}, ErrTokenExpired
			}
			if validErr.Errors&jwt.ValidationErrorMalformed > 0 {
				return CustomInfo{}, ErrTokenValid
			}
		}
		return CustomInfo{}, ErrTokenValid
	}
	if claims, ok := res.Claims.(*CustomClaims); ok && res.Valid {
		return CustomInfo{
			Id:       claims.Id,
			Nickname: claims.Nickname,
		}, nil
	}
	return CustomInfo{}, ErrTokenValid
}
