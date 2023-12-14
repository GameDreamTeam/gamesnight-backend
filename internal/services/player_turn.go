package services

import (
	"errors"
	"gamesnight/internal/database"
	"gamesnight/internal/models"
)

func (gs *GameService) GetPhraseToBeGuessed(currentPhrases models.PhraseStatusMap) (string, error) {

	if models.CurrentIndex >= len(currentPhrases.Phrases) {
		//Show EndTheGame
		return "", errors.New("index out of range")
	}

	// Get the phrase at the specified index
	phrase := currentPhrases.Phrases[models.CurrentIndex]

	return phrase.Input, nil
}

func (gs *GameService) HandlePlayerGuess(game models.Game, choice string) error {
	currentPhrases, err := gs.GetCurrentPhraseMap(game.GameId)
	if err != nil {
		return err
	}

	// Update the choice based on the request
	if choice == "guessed" {
		currentPhrases.Status[models.CurrentIndex] = models.Guessed
		(*game.Teams)[game.CurrentTeamIndex].Score += 10
		database.SetGame(&game)
	}

	gs.SetCurrentPhraseMap(game.GameId, currentPhrases)
	models.CurrentIndex += 1

	return nil
}
