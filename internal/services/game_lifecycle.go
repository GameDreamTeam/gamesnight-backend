package services

import (
	"fmt"
	"gamesnight/internal/database"
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

func (gs *GameService) CreateNewGame(playerId string) (*models.GameMeta, error) {
	gameId, err := GetKeyGenerator().CreateGameKey()

	if err != nil {
		fmt.Printf("Error in creating new game %s", err)
		return nil, err
	}

	existingGame, _ := database.GetGame(gameId)
	// if err != nil {
	// 	fmt.Printf("Error checking for existing game: %s", err)
	// 	return nil, err
	// }

	if existingGame != nil {
		fmt.Printf("Game with gameId %s already exists", gameId)
		// Handle this scenario, possibly by generating a new gameId or returning an error
		//Add a maximum recursion depth
		return gs.CreateNewGame(playerId)
	}

	gameMeta := models.GameMeta{
		GameId:    gameId,
		AdminId:   playerId,
		//update to get Current time function from utils
		CreatedAt: time.Now(),
		Players:   &[]models.Player{},
	}

	game := models.Game{
		GameId:    gameId,
		GameState: models.PlayersJoining,
	}

	// Use go routines here
	database.SetGameMeta(&gameMeta)
	database.SetGame(&game)

	return &gameMeta, nil
}

func (gs *GameService) JoinGame(gameId string, player *models.Player) (*models.GameMeta, error) {
	// Check the state of game here

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
		fmt.Println("Not able to set game")
		return nil, err
	}

	err = database.SetPlayerDetails(*player)
	if err != nil {
		fmt.Println("Not able to set player")
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

	//Minimum 2 players need to present otherwise it will throw out of bounds in array
	updatedGame := StartingCurrentAndNextPlayer(game)

	err = database.SetGame(updatedGame)
	if err != nil {
		return nil, err
	}

	return game, nil
}
