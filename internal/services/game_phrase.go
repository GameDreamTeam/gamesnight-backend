package services

import (
	"gamesnight/internal/database"
	"gamesnight/internal/models"
	"math/rand"
)

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

func (gs *GameService) GetGamePhrases(gameId string) (*models.PhraseList, error) {
	phrases, err := database.GetGamePhrases(gameId)
	if err != nil {
		return nil, err
	}

	return phrases, nil
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
