package services

import (
	"errors"
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
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

	// can add a validation for duplicate playerId

	if err != nil {
		// fmt.Printf("Error in creating player key %s", err)
		return nil, err
	}

	player := &models.Player{
		Id: &key,
	}
	return player, nil
}

func (ps *PlayerService) GetPlayerPhrases(playerId string) (*models.PhraseList, error) {
	// Fetch phrases for the player
	phrases, err := database.GetPlayerPhrases(playerId)
	if err != nil {
		return nil, err
	}

	return phrases, nil
}

func (ps *PlayerService) GetPlayerDetails(playerID string) (*models.Player, error) {
	return database.GetPlayerDetails(playerID)
}

func (ps *PlayerService) RemovePlayer(gameMeta *models.GameMeta, playerID string) (*models.GameMeta, error) {
	// Find the index of the player in the Players slice
	playerIndex := -1
	for i, player := range *gameMeta.Players {
		if *player.Id == playerID {
			playerIndex = i
			break
		}
	}

	// If the player is not found, return an error
	if playerIndex == -1 {
		return nil, errors.New("player not found in the game")
	}

	// Create a new slice excluding the player to be removed
	updatedPlayers := append((*gameMeta.Players)[:playerIndex], (*gameMeta.Players)[playerIndex+1:]...)

	// Update the gameMeta with the new slice
	gameMeta.Players = &updatedPlayers

	err := database.SetGameMeta(gameMeta)
	if err != nil {
		return nil, err
	}

	return gameMeta, nil
}

func (ps *PlayerService) NextPlayerAndTeam(gameId string) (*models.Game, error) {
	game, err := database.GetGame(gameId)
	if err != nil {
		return game, err
	}
	updateGame := ChangeNextPlayerAndTeam(game)

	err = database.SetGame(updateGame)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (ps *PlayerService) PlayerAlreadyAddedPhrases(playerId string) error {
	player, err := database.GetPlayerDetails(playerId)
	if err != nil {
		return err
	}
	if player.PhrasesSubmitted {
		logger.GetLogger().Logger.Error("player:" + playerId + " has already submitted phraes")
		return errors.New("you have already added phrases")
	}
	return nil
}
