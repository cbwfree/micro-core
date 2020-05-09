// JWT标准中注册的声明 (建议但不强制使用) ：
//
// iss: jwt签发者
// sub: jwt所面向的用户
// aud: 接收jwt的一方
// exp: jwt的过期时间，这个过期时间必须要大于签发时间
// nbf: 定义在什么时间之前，该jwt都是不可用的.
// iat: jwt的签发时间
// jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。
//
package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	opts *Options
}

func (t *Token) Encrypt(data interface{}) (string, error) {
	now := time.Now()

	// JWT声明
	claims := new(Claims)
	claims.Data = data
	claims.NotBefore = now.Unix()
	claims.Issuer = t.opts.Issuer
	claims.IssuedAt = now.Unix()

	if t.opts.Expire > 0 {
		claims.ExpiresAt = now.Add(t.opts.Expire).Unix()
	}

	token := jwt.NewWithClaims(t.opts.SigningMethod, claims)
	ts, err := token.SignedString(t.opts.SecretKey)
	if err != nil {
		return "", err
	}
	return ts, nil
}

func (t *Token) Verify(str string, result interface{}) (string, error) {
	claims := &Claims{Data: result}
	token, err := jwt.ParseWithClaims(str, claims, func(_ *jwt.Token) (interface{}, error) {
		return t.opts.SecretKey, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	// 检查是否需要刷新
	if t.opts.Refresh > 0 && time.Now().Add(-t.opts.Refresh).Unix() > claims.IssuedAt {
		return t.Encrypt(claims.Data)
	}

	return "", nil
}

func NewToken(opts ...Option) *Token {
	t := &Token{
		opts: new(Options),
	}
	t.opts.Init(opts...)
	return t
}
