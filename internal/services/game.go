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

	return gameMeta, nil
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

	//Minimum 2 players need to present otherwise it will throw out of bounds in array
	currentTeamIndex := game.CurrentTeamIndex
	nextTeamIndex := getNextTeamIndex(game.CurrentTeamIndex)
	currentTeamCurrentPlayerIndex := (*game.Teams)[currentTeamIndex].CurrentPlayerIndex
	nextTeamCurrentPlayerIndex := (*game.Teams)[nextTeamIndex].CurrentPlayerIndex

	game.CurrentPlayer = &(*(*game.Teams)[currentTeamIndex].Players)[currentTeamCurrentPlayerIndex]
	game.NextPlayer = &(*(*game.Teams)[nextTeamIndex].Players)[nextTeamCurrentPlayerIndex]

	// Randomize the phrase, make it in a map form
	phraseStatusMap, err := gs.GetRandomizedGamePhrases(gameId)
	if err != nil {
		return nil, err
	}

	// Store the PhraseStatusMap in Redis
	err = database.SetGamePhraseStatusMap(gameId, phraseStatusMap)
	if err != nil {
		return nil, err
	}

	return game, nil

}

func getNextTeamIndex(currentIndex int) int {
	if currentIndex == 1 {
		return 0
	}
	return 1
}

func (gs *GameService) RemovePlayer(gameMeta *models.GameMeta, playerID string) (*models.GameMeta, error) {
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
