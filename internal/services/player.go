package services

import (
	"fmt"
	"gamesnight/internal/models"

	"github.com/gin-gonic/gin"
)

type PlayerService struct{}

var us *PlayerService

func NewPlayerService() {
	us = &PlayerService{}
}

func GetPlayerService() *PlayerService {
	return us
}

func CreateNewPlayer() (*models.Player, error) {
	key, err := GetKeyGenerator().CreatePlayerKey()

	if err != nil {
		fmt.Printf("Error in creating player key %s", err)
		return nil, err
	}

	user := &models.Player{
		Id: &key,
	}
	return user, nil
}

func (us *PlayerService) GetPlayer(c *gin.Context, playerCookieName string) (*models.Player, error) {

	playerCookie, err := c.Cookie(playerCookieName)
	if err != nil {
		player, err := CreateNewPlayer()

		if err != nil {
			fmt.Printf("Error in creating new user %s", err)
			return nil, err
		}

		token, err := GetTokenService().CreatePlayerToken(*player.Id)
		if err != nil {
			fmt.Printf("Error in creating new token %s", err)
			return nil, err
		}
		// Need to check these other parameters
		c.SetCookie(playerCookieName, token.Token, 3600, "/", "", false, true)
		playerCookie = token.Token
	}

	player, err := GetTokenService().ParsePlayerToken(playerCookie)
	if err != nil {
		// Create new user in this case
		// Check if expired or not able to decode this token
		fmt.Printf("Error in parsing token %s", err)
		return nil, err
	}

	return player, nil
}
