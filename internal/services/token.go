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

func (ts *TokenService) CreateUserToken(userId string) (*models.Token, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	claims := &models.UserJWTClaims{
		UserId: userId,
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

func (ts *TokenService) ParseUserToken(tokenString string) (*models.User, error) {
	claims := &models.UserJWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		fmt.Println("Error in parsing token")
		return nil, err
	}

	// Check if the token is valid and has expected claims
	if claims, ok := token.Claims.(*models.UserJWTClaims); ok && token.Valid {
		return &models.User{
			UserId: &claims.UserId,
		}, nil
	} else {
		return nil, fmt.Errorf("Invalid token")
	}
}
