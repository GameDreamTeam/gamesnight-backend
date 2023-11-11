package services

import (
	"fmt"
	"gamesnight/internal/models"
	"time"
)

type GameService struct{}

var gs *GameService

func NewGameService() {
	gs = &GameService{}
}

func GetGameService() *GameService {
	return gs
}

func (gs *GameService) CreateNewGame(user *models.User) (*models.Game, error) {
	gameId, err := GetKeyGenerator().CreateGameKey()

	//Check if game already exists before returning this
	if err != nil {
		fmt.Printf("Error in creating new game %s", err)
		return nil, err
	}

	game := models.Game{
		GameId:    gameId,
		AdminId:   *user.UserId,
		CreatedAt: time.Now(),
	}

	return &game, nil
}

// func (gs *GameService) CreateNewGame(user *models.User) (*models.Game, error) {
