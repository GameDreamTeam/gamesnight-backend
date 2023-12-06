package services

import (
	"errors"
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
	// Check if game exists
	_, err := gs.GetGameMeta(gameId)
	if err != nil {
		return nil, err
	}

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

func (gs *GameService) GenerateRandom(phrases *models.PhraseList) (*models.PhraseStatusMap, error) {
	if phrases == nil || phrases.List == nil {
		// return error
		return &models.PhraseStatusMap{Phrases: make(map[string]models.PhraseStatus)}, nil
	}

	// Clone the original list to avoid modifying the original
	clonedList := make([]models.Phrase, len(*phrases.List))
	copy(clonedList, *phrases.List)

	// Use rand.Shuffle to randomize the list
	rand.Shuffle(len(clonedList), func(i, j int) {
		clonedList[i], clonedList[j] = clonedList[j], clonedList[i]
	})

	// Create a map with random phrases and empty string values
	phraseStatusMap := make(map[string]models.PhraseStatus)
	for _, phrase := range clonedList {
		phraseStatusMap[phrase.Input] = models.NotGuessed
	}

	// Return the new PhraseStatusMap with randomized phrases
	return &models.PhraseStatusMap{Phrases: phraseStatusMap}, nil
}

func (gs *GameService) SetCurrentPhrases(gameId string, currentPhrases *models.PhraseStatusMap) error {
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
	err = database.SetCurrentPhrases(gameId, *currentPhrases)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GameService) GetCurrentPhrases(gameId string) (models.PhraseStatusMap, error) {
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
	currentPhrases, err := database.GetCurrentPhrases(gameId)
	if err != nil {
		return models.PhraseStatusMap{}, err
	}

	return currentPhrases, nil
}

func (gs *GameService) GetNextPhrase(currentPhrases models.PhraseStatusMap, index int) (string, error) {
	var keys []string
	for key := range currentPhrases.Phrases {
		keys = append(keys, key)
	}

	// Check if the index is within range
	if index < 0 || index >= len(keys) {
		return "", errors.New("Index out of range")
	}

	// Get the phrase at the specified index
	phrase := keys[index]

	return phrase, nil
}

func (gs *GameService) HandlePlayerGuess(gameId string, playerId *string, choice string, key string) error {
	currentPhrases, err := gs.GetCurrentPhrases(gameId)
	if err != nil {
		return err
	}

	// Update the choice based on the request
	if choice == "guessed" {
		currentPhrases.Phrases[key] = models.Guessed
	} else {
		// no change needed
	}

	gs.SetCurrentPhrases(gameId, &currentPhrases)

	return nil
}
