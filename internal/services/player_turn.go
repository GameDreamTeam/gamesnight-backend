package services

import (
	"errors"
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"

	"go.uber.org/zap"
)

func (gs *GameService) GetPhraseToBeGuessed(currentPhrases models.PhraseStatusMap, game models.Game) (string, error) {
	phraseIndex := game.CurrentPhraseMapIndex
	if phraseIndex >= len(currentPhrases.Phrases) {
		game.GameState = models.Finished
		database.SetGame(&game)
		return "the game has ended", nil
	}

	phrase := currentPhrases.Phrases[phraseIndex]
	return phrase.Input, nil
}

func (gs *GameService) HandlePlayerGuess(game models.Game, choice string) (models.PhraseStatusMap, error) {
	currentPhraseMap, err := gs.GetCurrentPhraseMap(game)
	if err != nil {
		return currentPhraseMap, err
	}

	if choice == "guessed" {
		currentPhraseMap.Status[game.CurrentPhraseMapIndex] = models.Guessed
		(*game.Teams)[game.CurrentTeamIndex].Score += 10
	}
	game.CurrentPhraseMapIndex += 1
	database.SetGame(&game)
	gs.SetCurrentPhraseMap(game.GameId, currentPhraseMap)

	return currentPhraseMap, nil
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
