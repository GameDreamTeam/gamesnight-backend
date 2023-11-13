package services

import (
	"errors"
	"fmt"
	"gamesnight/internal/database"
	"gamesnight/internal/models"
	"math/rand"
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

	gameMeta := models.GameMeta{
		GameId:    gameId,
		AdminId:   playerId,
		CreatedAt: time.Now(),
		Players:   &[]models.Player{},
	}

	database.SetGameMeta(&gameMeta)

	return &gameMeta, nil
}

func (gs *GameService) JoinGame(gameId string, player *models.Player) (*models.GameMeta, error) {

	// This entire portion has to acquire a lock when having high concurrency
	game, err := database.GetGameMeta(gameId)
	if err != nil {
		return nil, err
	}

	game, err = addPlayerToGame(game, player)
	if err != nil {
		return nil, err
	}

	err = database.SetGameMeta(game)

	if err != nil {
		fmt.Println("Not able to set game")
		return nil, err
	}

	return game, nil
}

func (gs *GameService) GetGameMeta(gameId string) (*models.GameMeta, error) {
	return database.GetGameMeta(gameId)
}

func (gs *GameService) StartGame(gamemeta *models.GameMeta) (*models.Game, error) {
	// Future we have to make number of teams customizable
	team1, team2 := dividePlayersIntoTeams(*gamemeta.Players)

	// Make these names customizable
	t1 := models.Team{
		Name:    "RED",
		Players: &team1,
	}

	t2 := models.Team{
		Name:    "BLUE",
		Players: &team2,
	}

	teams := []models.Team{t1, t2}

	game := models.Game{
		GameId:    gamemeta.GameId,
		Teams:     teams,
		GameState: models.Playing,
	}

	return &game, nil
}

func dividePlayersIntoTeams(players []models.Player) ([]models.Player, []models.Player) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	mid := len(players) / 2
	return players[:mid], players[mid:]
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
