package services

import (
	"errors"
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

	//Check if game already exists before returning this
	if err != nil {
		fmt.Printf("Error in creating new game %s", err)
		return nil, err
	}

	game := models.GameMeta{
		GameId:    gameId,
		AdminId:   playerId,
		CreatedAt: time.Now(),
		Players:   &[]models.Player{},
	}

	database.SetGame(&game)

	return &game, nil
}

func (gs *GameService) JoinGame(gameId string, player *models.Player) (*models.GameMeta, error) {

	// This entire portion has to acquire a lock when having high concurrency
	game, err := database.GetGame(gameId)
	if err != nil {
		return nil, err
	}

	game, err = addPlayerToGame(game, player)
	if err != nil {
		return nil, err
	}

	err = database.SetGame(game)

	if err != nil {
		fmt.Println("Not able to set game")
		return nil, err
	}

	return game, nil
}

func (gs *GameService) GetGame(gameId string) (*models.GameMeta, error) {
	return database.GetGame(gameId)
}

func addPlayerToGame(game *models.GameMeta, player *models.Player) (*models.GameMeta, error) {

	if !contains(*game.Players, player) {
		*game.Players = append(*game.Players, *player)
	} else {
		// Return custom error here (404)
		return nil, errors.New("player already exists in this game")
	}

	return game, nil
}

func contains(playerSlice []models.Player, player *models.Player) bool {
	for _, p := range playerSlice {
		if *p.Id == *player.Id {
			return true
		}
	}
	return false
}
