package services

import (
	"gamesnight/internal/database"
	"gamesnight/internal/models"
)

func (gs *GameService) StartGame(gameId string) (*models.Game, error) {
	//Check if game is teams divided and ready to start

	game, err := database.GetGame(gameId)
	if err != nil {
		return nil, err
	}

	// Need to check the current status of game before starting game

	//Minimum 2 players need to present otherwise it will throw out of bounds in array
	// Name of method should be a verb
	updatedGame := StartingCurrentAndNextPlayer(game)

	err = database.SetGame(updatedGame)
	if err != nil {
		return nil, err
	}

	return game, nil
}
