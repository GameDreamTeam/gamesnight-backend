package models

import "github.com/golang-jwt/jwt"

type PlayerJWTClaims struct {
	PlayerId string `json:"username"`
	jwt.StandardClaims
}

type Token struct {
	Token string `json:"token"`
}
