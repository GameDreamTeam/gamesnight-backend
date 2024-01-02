package services

import (
	"gamesnight/internal/database"
	"gamesnight/internal/models"
)

func (gs *GameService) MakeTeams(gamemeta *models.GameMeta) (*models.Game, error) {
	//Check if game already exists or not before making teams

	// Need to acquire a lock before setting this team
	game, err := database.GetGame(gamemeta.GameId)
	if err != nil {
		return nil, err
	}
	// If we have fetched the game then we should check the status of game.
	//If it is already teams divided or further ahead then we should not allow this

	// Check if atleast 2 players exist in the game
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
