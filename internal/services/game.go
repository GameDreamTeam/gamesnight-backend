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

func (gs *GameService) MakeTeams(gamemeta *models.GameMeta) (*models.Game, error) {
	//Check if game already exists or not before making teams

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
		Teams:     &teams,
		GameState: models.Playing,
	}

	// Write this to redis

	return &game, nil
}

func (gs *GameService) StartGame(gameId string) (*models.Game, error) {
	//Check if game is teams divided and ready to start

	game, err := database.GetGame(gameId)
	if err != nil {
		return nil, err
	}

	game.GameState = models.Playing

	err = database.SetGame(game)
	if err != nil {
		return nil, err
	}

	return game, nil

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
