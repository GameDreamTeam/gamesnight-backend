package services

import (
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"
)

type GameService struct{}

var gs *GameService

func NewGameService() {
	gs = &GameService{}
}

func GetGameService() *GameService {
	return gs
}

func (gs *GameService) CreateNewGame(playerId string) (*models.GameMeta, error) {
	gameId, err := GetKeyGenerator().CreateGameKey()

	if err != nil {
		return nil, err
	}

	existingGame, _ := database.GetGame(gameId)

	if existingGame != nil {
		return gs.CreateNewGame(playerId)
	}

	// Logic of creating two entities is that one can stay mostly constant and only other is varying
	// Also we should see what information would be fetched often. For that information we should
	// reduce the number of network calls we make to redis, hence we can store most information in same
	// object of Game and GameMeta

	gameMeta := models.GameMeta{
		GameId:    gameId,
		AdminId:   playerId,
		CreatedAt: GetCurrentTime(),
		Players:   &[]models.Player{},
	}

	game := models.Game{
		GameId:    gameId,
		GameState: models.PlayersJoining,
	}

	logger.GetLogger().Logger.Info("player:" + playerId + " created game:" + gameId + " successfully")

	// Use go routines here for concurrency and better speed
	database.SetGameMeta(&gameMeta)
	database.SetGame(&game)

	return &gameMeta, nil
}

func (gs *GameService) JoinGame(gameMeta *models.GameMeta, player *models.Player) (*models.GameMeta, error) {
	// This entire portion has to acquire a lock when having high concurrency

	gameMetaWithPlayer, err := addPlayerToGame(gameMeta, player)
	if err != nil {
		return nil, err
	}

	err = database.SetGameMeta(gameMetaWithPlayer)
	if err != nil {
		return nil, err
	}

	logger.GetLogger().Logger.Info("player:" + *player.Id + " joined game:" + gameMeta.GameId + " successfully")

	err = database.SetPlayerDetails(*player)
	if err != nil {
		return nil, err
	}

	return gameMeta, nil
}

func (gs *GameService) UpdateStateOfGame(gameId string) (*models.Game, error) {
	game, err := database.GetGame(gameId)
	if err != nil {
		return nil, err
	}

	game.GameState += 1

	err = database.SetGame(game)
	if err != nil {
		return nil, err
	}

	currentGameState := GetGameState(game.GameState)
	logger.GetLogger().Logger.Info("GameState of game:" + gameId + " updated to " + currentGameState)

	return game, nil
}
