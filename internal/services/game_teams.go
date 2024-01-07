package services

import (
	"errors"
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"

	"go.uber.org/zap"
)

func (gs *GameService) MakeTeams(gamemeta *models.GameMeta) (*models.Game, error) {
	// Need to acquire a lock before setting this team
	game, err := database.GetGame(gamemeta.GameId)
	if err != nil {
		return nil, err
	}

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

	err = database.SetGame(game)
	if err != nil {
		return nil, err
	}

	logger.GetLogger().Logger.Info(
		"Teams for game:"+game.GameId+" divided successfull",
		zap.Any("teams", game.Teams))

	return game, nil
}

func (gs *GameService) CheckIfAllPlayerHaveSubmittedPhrases(gameMeta models.GameMeta) error {
	for _, player := range *gameMeta.Players {
		if !player.PhrasesSubmitted {
			return errors.New("all players have not submitted words")
		}
	}
	return nil
}
