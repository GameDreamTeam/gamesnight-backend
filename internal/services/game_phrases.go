package services

import (
	"errors"
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"
	"math/rand"

	"go.uber.org/zap"
)

func (gs *GameService) PlayerExistInGame(gameMeta models.GameMeta, player models.Player) error {
	if !contains(*gameMeta.Players, &player) {
		logger.GetLogger().Logger.Error("player:" + *player.Id + " has not joined game:" + gameMeta.GameId)
		return errors.New("player not found in the game")
	}
	return nil
}

func (gs *GameService) AddPhrasesToGame(playerId string, gameMeta *models.GameMeta, phraseList *models.PhraseList) error {
	newGameMeta, err := MarkPlayerHasAddedWords(gameMeta, playerId)
	if err != nil {
		return err
	}

	err = database.SetGameMeta(&newGameMeta)
	if err != nil {
		return err
	}

	err = database.SetGamePhrases(gameMeta.GameId, phraseList)
	if err != nil {
		return err
	}

	logger.GetLogger().Logger.Info(
		"player:"+playerId+" in game:"+gameMeta.GameId+"added words successfully",
		zap.Any("phraseList", phraseList),
	)

	return nil
}

func (gs *GameService) AddPhrasesToPlayer(player models.Player, phraseList *models.PhraseList) error {
	err := database.SetPlayerPhrases(*player.Id, phraseList)
	if err != nil {
		return err
	}

	player.PhrasesSubmitted = true
	err = database.SetPlayerDetails(player)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GameService) SetCurrentPhraseMap(gameId string, currentPhrases models.PhraseStatusMap) error {
	return database.SetCurrentPhraseMap(gameId, currentPhrases)
}

func (gs *GameService) GetCurrentPhraseMap(game models.Game) (models.PhraseStatusMap, error) {
	currentGamePhraseMap, err := database.GetCurrentPhraseMap(game.GameId)
	if err != nil {
		return currentGamePhraseMap, err
	}
	if len(currentGamePhraseMap.Phrases) == 0 {
		game.GameState = models.Finished
		database.SetGame(&game)
		return currentGamePhraseMap, errors.New("all phrase are completed")
	}
	return currentGamePhraseMap, nil
}

func (gs *GameService) RandomizePhrases(phraseMap models.PhraseStatusMap) models.PhraseStatusMap {
	rand.Shuffle(len(phraseMap.Phrases), func(i, j int) {
		phraseMap.Phrases[i], phraseMap.Phrases[j] = phraseMap.Phrases[j], phraseMap.Phrases[i]
		phraseMap.Status[i], phraseMap.Status[j] = phraseMap.Status[j], phraseMap.Status[i]
	})
	return phraseMap
}

func (gs *GameService) RemoveGuessedPhrases(gameId string, phraseMap models.PhraseStatusMap) models.PhraseStatusMap {
	var newPhrases []models.Phrase
	var newStatus []models.PhraseStatus

	for i, status := range phraseMap.Status {
		if status != models.Guessed {
			newPhrases = append(newPhrases, phraseMap.Phrases[i])
			newStatus = append(newStatus, status)
		}
	}

	newMap := models.PhraseStatusMap{
		Phrases: newPhrases,
		Status:  newStatus,
	}

	randomizedMap := gs.RandomizePhrases(newMap)

	database.SetCurrentPhraseMap(gameId, randomizedMap)

	return newMap
}
