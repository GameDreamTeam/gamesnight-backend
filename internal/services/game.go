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
		PlayerIds: &[]string{*user.UserId},
	}

	database.SetGame(&game)

	return &game, nil
}

func (gs *GameService) JoinGame(gameId string, userId string) (*models.Game, error) {

	game, err := database.GetGame(gameId)
	if err != nil {
		fmt.Println("Not getting game")
		return nil, err
	}

	game, err = addPlayer(game, userId)
	if err != nil {
		fmt.Println("Not able to add player")
		return nil, err
	}

	err = database.SetGame(game)

	if err != nil {
		fmt.Println("Not able to set game")
		return nil, err
	}

	return game, nil
}

func (gs *GameService) GetGame(gameId string) (*models.Game, error) {
	return database.GetGame(gameId)
}

func addPlayer(game *models.Game, userId string) (*models.Game, error) {
	*game.PlayerIds = append(*game.PlayerIds, userId)
	return game, nil
}
