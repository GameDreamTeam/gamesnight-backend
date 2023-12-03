package services

import (
	"fmt"
	"gamesnight/internal/database"
	"gamesnight/internal/models"
)

type PlayerService struct{}

var ps *PlayerService

func NewPlayerService() {
	ps = &PlayerService{}
}

func GetPlayerService() *PlayerService {
	return ps
}

func (ps *PlayerService) CreateNewPlayer() (*models.Player, error) {
	key, err := GetKeyGenerator().CreatePlayerKey()

	if err != nil {
		fmt.Printf("Error in creating player key %s", err)
		return nil, err
	}

	player := &models.Player{
		Id: &key,
	}
	return player, nil
}

func (ps *PlayerService) GetPlayerDetails(playerID string) (*models.Player, error) {
	return database.GetPlayerDetails(playerID)
}
