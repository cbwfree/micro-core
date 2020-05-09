package jwt

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	jwt.StandardClaims
	Data interface{} `json:"data,omitempty"`
}
