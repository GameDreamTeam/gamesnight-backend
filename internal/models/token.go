package models

import "github.com/golang-jwt/jwt"

type UserJWTClaims struct {
	UserId string `json:"username"`
	jwt.StandardClaims
}

type Token struct {
	Token string `json:"token"`
}
