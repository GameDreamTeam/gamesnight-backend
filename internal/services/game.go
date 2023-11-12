package services

import (
	"fmt"
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"
	"time"

	"go.uber.org/zap"
)

type GameService struct{}

var gs *GameService

func NewGameService() {
	gs = &GameService{}
}

func GetGameService() *GameService {
	return gs
}

func (gs *GameService) CreateNewGame(player *models.Player) (*models.Game, error) {
	gameId, err := GetKeyGenerator().CreateGameKey()

	//Check if game already exists before returning this
	if err != nil {
		fmt.Printf("Error in creating new game %s", err)
		return nil, err
	}

	game := models.Game{
		GameId:    gameId,
		Admin:     player,
		CreatedAt: time.Now(),
		PlayerIds: &[]models.Player{},
	}

	database.SetGame(&game)

	return &game, nil
}

func (gs *GameService) JoinGame(gameId string, player *models.Player) (*models.Game, error) {

	// This entire portion has to acquire a lock when having high concurrency
	game, err := database.GetGame(gameId)
	if err != nil {
		fmt.Println("Not getting game")
		return nil, err
	}

	game, err = addPlayer(game, player)
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

func addPlayer(game *models.Game, player *models.Player) (*models.Game, error) {

	logger.GetLogger().Logger.Info("Player add", zap.Any("game", *game), zap.Any("player", *player))
	if !contains(*game.PlayerIds, player) {
		*game.PlayerIds = append(*game.PlayerIds, *player)
	}

	return game, nil
}

// Move this to util class
func contains(playerSlice []models.Player, player *models.Player) bool {
	for _, p := range playerSlice {
		if *p.Id == *player.Id {
			return true
		}
	}
	return false
}
