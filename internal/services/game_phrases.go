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

func (gs *GameService) GetRandomizedGamePhrases(gameID string) (models.PhraseStatusMap, error) {
	phraseList, err := gs.GetGamePhrases(gameID)
	if err != nil {
		return models.PhraseStatusMap{}, err
	}

	// Return an empty map if no phrases are available
	if len(*phraseList.List) == 0 {
		return models.PhraseStatusMap{}, errors.New("no phrases found for the game")
	}

	// Randomize the order of phrases
	randomizedPhrases := models.PhraseStatusMap{Phrases: make(map[string]models.PhraseStatus)}
	for _, phrase := range *phraseList.List {
		randomizedPhrases.Phrases[phrase.Input] = models.NotGuessed // Using Input as the key
	}

	rand.Shuffle(len(*phraseList.List), func(i, j int) {
		phrases := *phraseList.List
		randomizedPhrases.Phrases[phrases[i].Input], randomizedPhrases.Phrases[phrases[j].Input] = randomizedPhrases.Phrases[phrases[j].Input], randomizedPhrases.Phrases[phrases[i].Input]
	})

	return randomizedPhrases, nil
}
