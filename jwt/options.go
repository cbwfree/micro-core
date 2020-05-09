package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Option func(*Options)

type Options struct {
	SigningMethod jwt.SigningMethod
	SecretKey     []byte
	Issuer        string
	Expire        time.Duration
	Refresh       time.Duration
}

func (o *Options) Init(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// 签名方式
func SigningMethod(method jwt.SigningMethod) Option {
	return func(o *Options) {
		o.SigningMethod = method
	}
}

// 秘钥信息
func SecretKey(secret string) Option {
	return func(o *Options) {
		o.SecretKey = []byte(secret)
	}
}

// 签发者
func Issuer(issuer string) Option {
	return func(o *Options) {
		o.Issuer = issuer
	}
}

// 有效期
func Expire(expire time.Duration) Option {
	return func(o *Options) {
		o.Expire = expire
	}
}

// 刷新间隔
func Refresh(refresh time.Duration) Option {
	return func(o *Options) {
		o.Refresh = refresh
	}
}
