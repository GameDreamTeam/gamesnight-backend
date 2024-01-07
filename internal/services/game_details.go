package services

import (
	"gamesnight/internal/database"
	"gamesnight/internal/models"
)

func (gs *GameService) GetGameMeta(gameId string) (*models.GameMeta, error) {
	return database.GetGameMeta(gameId)
}

func (gs *GameService) GetGame(gameId string) (*models.Game, error) {
	return database.GetGame(gameId)
}

func (gs *GameService) GetGamePhrases(gameId string) (*models.PhraseList, error) {
	return database.GetGamePhrases(gameId)
}
