package models

import "github.com/golang-jwt/jwt"

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
