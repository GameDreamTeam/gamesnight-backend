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
	err := database.SetCurrentPhraseMap(gameId, currentPhrases)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GameService) GetCurrentPhraseMap(gameId string) (models.PhraseStatusMap, error) {
	currentPhrases, err := database.GetCurrentPhraseMap(gameId)
	if err != nil {
		return models.PhraseStatusMap{}, err
	}

	return currentPhrases, nil
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

	rand.Shuffle(len(newPhrases), func(i, j int) { newPhrases[i], newPhrases[j] = newPhrases[j], newPhrases[i] })

	newMap := models.PhraseStatusMap{
		Phrases: newPhrases,
		Status:  newStatus,
	}
	database.SetCurrentPhraseMap(gameId, newMap)

	return newMap
}
