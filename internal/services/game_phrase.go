package services

import (
	"gamesnight/internal/database"
	"gamesnight/internal/models"
	"math/rand"
)

func (gs *GameService) AddPhrasesToGame(playerId string, gameId string, phraseList *models.PhraseList) error {
	game, err := database.GetGame(gameId)
	if err != nil {
		return err
	}

	if game.GameState != models.AddingWords {
		if err != nil {
			return err
		}
	}
	// Add phrases to the game
	gameMeta, err := gs.GetGameMeta(gameId)

	newGame, err := MarkPlayerHasAddedWords(gameMeta, playerId)

	if err != nil {
		return err
	}
	err = database.SetGame(game)
	err = database.SetGameMeta(&newGame)

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

func (gs *GameService) StartGame(gameId string) (*models.Game, error) {
	//Check if game is teams divided and ready to start

	game, err := database.GetGame(gameId)
	if err != nil {
		return nil, err
	}

	// Need to check the current status of game before starting game

	//Minimum 2 players need to present otherwise it will throw out of bounds in array
	// Name of method should be a verb
	updatedGame := StartingCurrentAndNextPlayer(game)

	err = database.SetGame(updatedGame)
	if err != nil {
		return nil, err
	}

	return game, nil
}
