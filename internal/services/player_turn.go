package services

import (
	"errors"
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"
)

func (gs *GameService) GetPhraseToBeGuessed(currentPhrases models.PhraseStatusMap, phraseIndex int) (string, error) {

	if phraseIndex >= len(currentPhrases.Phrases) {
		//Show EndTheGame
		return "", errors.New("index out of range")
	}

	// Get the phrase at the specified index
	phrase := currentPhrases.Phrases[phraseIndex]

	return phrase.Input, nil
}

func (gs *GameService) HandlePlayerGuess(game models.Game, choice string) error {
	currentPhrases, err := gs.GetCurrentPhraseMap(game.GameId)
	if err != nil {
		return err
	}

	// Update the choice based on the request
	if choice == "guessed" {
		currentPhrases.Status[game.CurrentPhraseMapIndex] = models.Guessed
		(*game.Teams)[game.CurrentTeamIndex].Score += 10
		game.CurrentPhraseMapIndex += 1
		database.SetGame(&game)
	}

	gs.SetCurrentPhraseMap(game.GameId, currentPhrases)

	return nil
}

func (gs *GameService) CheckCurrentPlayer(gameId string, playerId string) error {
	if playerId != gameId {
		logger.GetLogger().Logger.Error(
			"player starting turn should be current player",
		)
		return errors.New("you are not the current player")

	}
	return nil
}
