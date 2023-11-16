package models

import "github.com/golang-jwt/jwt"

type PlayerJWTClaims struct {
	PlayerId string `json:"playerId"`
	jwt.StandardClaims
}

type Token struct {
	Token string `json:"token"`
}
