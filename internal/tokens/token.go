package tokens

import (
	"gamesnight/internal/models"
	"time"

	"github.com/golang-jwt/jwt"
)

type Tokens struct{}

var jwtKey = []byte("games-secret-for-us")
var tokens *Tokens

func NewToken() {
	tokens = &Tokens{}
}

func GetTokens() *Tokens {
	return tokens
}

func CreateUserToken(uname *models.User) (*string, error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &models.JWTClaims{
		Username: uname.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
