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

func (gs *GameService) GetGameMeta(gameId string) (*models.GameMeta, error) {
	return database.GetGameMeta(gameId)
}

func (gs *GameService) GetGame(gameId string) (*models.Game, error) {
	return database.GetGame(gameId)
}

// func (gs *GameService) StartTurn(player *models.Player) (*models.Game, error) {

// }

func (gs *GameService) MakeTeams(gamemeta *models.GameMeta) (*models.Game, error) {
	//Check if game already exists or not before making teams

	// Need to acquire a lock before setting this team
	game, err := database.GetGame(gamemeta.GameId)
	if err != nil {
		return nil, err
	}

	// Future we have to make number of teams customizable
	team1, team2 := dividePlayersIntoTeams(*gamemeta.Players)

	// Make these names customizable
	t1 := models.Team{
		Name:               "RED",
		Players:            &team1,
		Score:              0,
		CurrentPlayerIndex: 0,
	}

	t2 := models.Team{
		Name:               "BLUE",
		Players:            &team2,
		Score:              0,
		CurrentPlayerIndex: 0,
	}

	teams := []models.Team{t1, t2}
	game.Teams = &teams
	game.GameState = models.TeamsDivided

	// Write this to redis
	err = database.SetGame(game)
	if err != nil {
		return nil, err
	}

	return game, nil
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
	return game, nil

}

func dividePlayersIntoTeams(players []models.Player) ([]models.Player, []models.Player) {
	// if team exits in
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

func getNextTeamIndex(currentIndex int) int {
	if currentIndex == 1 {
		return 0
	}
	return 1
}

func (gs *GameService) AddPhrasesToGame(gameId string, phraseList *models.PhraseList) error {
	// Check if game exists
	game, err := gs.GetGame(gameId)
	if err != nil {
		return err
	}

	if game.GameState != models.AddingWords {
		game.GameState = models.AddingWords
		err = database.SetGame(game)
		if err != nil {
			return err
		}
	}
	// Add phrases to the game
	err = database.SetGamePhrases(gameId, phraseList)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GameService) AddPhrasesToPlayer(playerId string, phraseList *models.PhraseList) error {
	// Add phrases to the player
	err := database.SetPlayerPhrases(playerId, phraseList)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GameService) GetGamePhrases(gameId string) (*models.PhraseList, error) {
	// Check if game exists
	_, err := gs.GetGameMeta(gameId)
	if err != nil {
		return nil, err
	}

	// Fetch phrases for the game
	phrases, err := database.GetGamePhrases(gameId)
	if err != nil {
		return nil, err
	}

	return phrases, nil
}

func (ps *PlayerService) GetPlayerPhrases(playerId string) (*models.PhraseList, error) {
	// Fetch phrases for the player
	phrases, err := database.GetPlayerPhrases(playerId)
	if err != nil {
		return nil, err
	}

	return phrases, nil
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
		return nil, errors.New("Player not found in the game")
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
