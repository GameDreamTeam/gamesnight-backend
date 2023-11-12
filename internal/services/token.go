package services

import (
	"fmt"
	"gamesnight/internal/models"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService struct{}

var ts *TokenService

var jwtKey = []byte("games-secret-for-us")

func NewTokenService() {
	ts = &TokenService{}
}

func GetTokenService() *TokenService {
	return ts
}

func (ts *TokenService) CreatePlayerToken(playerId string) (*models.Token, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	claims := &models.PlayerJWTClaims{
		PlayerId: playerId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return nil, err
	}

	t := &models.Token{
		Token: tokenString,
	}

	return t, nil
}

func (ts *TokenService) ParsePlayerToken(tokenString string) (*models.Player, error) {
	claims := &models.PlayerJWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		fmt.Println("Error in parsing token")
		return nil, err
	}

	// Check if the token is valid and has expected claims
	if claims, ok := token.Claims.(*models.PlayerJWTClaims); ok && token.Valid {
		return &models.Player{
			Id: &claims.PlayerId,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
