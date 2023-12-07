package services

import (
	"errors"
	"gamesnight/internal/database"
	"gamesnight/internal/models"
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

func (gs *GameService) AddPhrasesToPlayer(playerId string, phraseList *models.PhraseList) error {
	// Add validation for playerId exists or not in the current game
	// Add phrases to the player
	err := database.SetPlayerPhrases(playerId, phraseList)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GameService) GetGamePhrases(gameId string) (*models.PhraseList, error) {
	// Check if game exists, no need for it
	// _, err := gs.GetGameMeta(gameId)
	// if err != nil {
	// 	return nil, err
	// }

	// Fetch phrases for the game
	phrases, err := database.GetGamePhrases(gameId)
	if err != nil {
		return nil, err
	}

	return phrases, nil
}

func (ps *PlayerService) GetPlayerPhrases(playerId string) (*models.PhraseList, error) {
	// Fetch phrases for the player
	phrases, err := database.GetPlayerPhrases(playerId)
	if err != nil {
		return nil, err
	}

	return phrases, nil
}

func (gs *GameService) GeneratePhraseListToMap(phrases *models.PhraseList) (models.PhraseStatusMap, error) {
	if phrases == nil || phrases.List == nil {
		// return error
		return models.PhraseStatusMap{}, errors.New("no phrases found")
	}

	// Clone the original list to avoid modifying the original
	// clonedList := make([]models.Phrase, len(*phrases.List))
	// copy(clonedList, *phrases.List)

	// Use rand.Shuffle to randomize the list
	// rand.Shuffle(len(clonedList), func(i, j int) {
	// 	clonedList[i], clonedList[j] = clonedList[j], clonedList[i]
	// })

	// Create a map with random phrases and empty string values
	phraseStatusMap := models.PhraseStatusMap{
		Phrases: make([]models.Phrase, len(*phrases.List)),
		Status:  make([]models.PhraseStatus, len(*phrases.List)),
	}

	for i, phrase := range *phrases.List {
		phraseStatusMap.Phrases[i] = phrase
		phraseStatusMap.Status[i] = models.NotGuessed
	}

	// Return the new PhraseStatusMap with randomized phrases
	return phraseStatusMap, nil
}

func (gs *GameService) SetCurrentPhraseMap(gameId string, currentPhrases models.PhraseStatusMap) error {
	// Check if game exists
	// game, err := gs.GetGame(gameId)
	// if err != nil {
	// 	return err
	// }

	// if game.GameState != models.AddingWords {
	// 	game.GameState = models.AddingWords
	// 	err = database.SetGame(game)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// Add phraseMap to the game
	err := database.SetCurrentPhraseMap(gameId, currentPhrases)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GameService) GetCurrentPhraseMap(gameId string) (models.PhraseStatusMap, error) {
	// Check if game exists
	game, err := gs.GetGame(gameId)
	if err != nil {
		return models.PhraseStatusMap{}, err
	}

	if game.GameState != models.Playing {
		game.GameState = models.Playing
		err = database.SetGame(game)
		if err != nil {
			return models.PhraseStatusMap{}, err
		}
	}

	// Get the current phrase map from redis
	currentPhrases, err := database.GetCurrentPhraseMap(gameId)
	if err != nil {
		return models.PhraseStatusMap{}, err
	}

	return currentPhrases, nil
}

// Fetches Phrase at index i from PhraseMap
func (gs *GameService) GetPhraseToBeGuessed(currentPhrases models.PhraseStatusMap) (string, error) {

	// Check if the index is within range
	//Check if Phrase is Guessed or not
	if models.CurrentIndex >= len(currentPhrases.Phrases) {
		return "", errors.New("index out of range")
	}

	// Get the phrase at the specified index
	phrase := currentPhrases.Phrases[models.CurrentIndex]

	return phrase.Input, nil
}

func (gs *GameService) HandlePlayerGuess(gameId string, playerId *string, choice string) error {
	currentPhrases, err := gs.GetCurrentPhraseMap(gameId)
	if err != nil {
		return err
	}

	// Update the choice based on the request
	if choice == "guessed" {
		currentPhrases.Status[models.CurrentIndex] = models.Guessed
	}

	gs.SetCurrentPhraseMap(gameId, currentPhrases)
	models.CurrentIndex += 1

	return nil
}
