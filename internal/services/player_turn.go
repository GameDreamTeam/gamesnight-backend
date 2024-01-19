package services

import (
	"errors"
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"

	"go.uber.org/zap"
)

func (gs *GameService) GetPhraseToBeGuessed(currentPhrases models.PhraseStatusMap, phraseIndex int) (string, error) {

	if phraseIndex >= len(currentPhrases.Phrases) {
		//Show EndTheGame
		return "the game has ended", errors.New("index out of range")
	}

	phrase := currentPhrases.Phrases[phraseIndex]

	return phrase.Input, nil
}

func (gs *GameService) HandlePlayerGuess(game models.Game, choice string) error {
	currentPhrases, err := gs.GetCurrentPhraseMap(game.GameId)
	if err != nil {
		return err
	}

	if choice == "guessed" {
		currentPhrases.Status[game.CurrentPhraseMapIndex] = models.Guessed
		(*game.Teams)[game.CurrentTeamIndex].Score += 10
		game.CurrentPhraseMapIndex += 1
		database.SetGame(&game)
	}

	gs.SetCurrentPhraseMap(game.GameId, currentPhrases)

	return nil
}

func (gs *GameService) CheckCurrentPlayer(playerId string, gameCurrentPlayer string) error {
	if playerId != gameCurrentPlayer {
		logger.GetLogger().Logger.Error(
			"player starting turn should be current player",
			zap.Any("player", playerId),
			zap.Any("gameCurrent", gameCurrentPlayer),
		)
		return errors.New("you are not the current player")

	}
	return nil
}
