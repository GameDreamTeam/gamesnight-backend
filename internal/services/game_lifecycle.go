package services

import (
	"gamesnight/internal/database"
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

	// Add log here that player xyz created game abc
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

	// Use go routines here for concurrency and better speed
	database.SetGameMeta(&gameMeta)
	database.SetGame(&game)

	return &gameMeta, nil
}

func (gs *GameService) JoinGame(gameId string, player *models.Player) (*models.GameMeta, error) {
	// This entire portion has to acquire a lock when having high concurrency
	gameMeta, err := database.GetGameMeta(gameId)
	if err != nil {
		return nil, err
	}

	gameMeta, err = addPlayerToGame(gameMeta, player)
	if err != nil {
		return nil, err
	}

	err = database.SetGameMeta(gameMeta)
	if err != nil {
		return nil, err
	}

	err = database.SetPlayerDetails(*player)
	if err != nil {
		return nil, err
	}

	return gameMeta, nil
}

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
