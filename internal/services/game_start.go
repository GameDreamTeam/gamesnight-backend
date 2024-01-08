package services

import (
	"gamesnight/internal/database"
	"gamesnight/internal/models"
)

func (gs *GameService) StartGame(game *models.Game) (*models.Game, error) {
	updatedGame := StartingCurrentAndNextPlayer(game)

	err := database.SetGame(updatedGame)
	if err != nil {
		return nil, err
	}

	return updatedGame, nil
}
