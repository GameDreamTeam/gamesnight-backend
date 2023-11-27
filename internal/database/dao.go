package database

import (
	"encoding/json"
	"fmt"
	"gamesnight/internal/models"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func SetGame(game *models.Game) error {
	key := GetGameKey(game.GameId)

	jsonGame, err := json.Marshal(game)
	if err != nil {
		return errors.Wrap(err, "Game json conversion failed while setting game")
	}

	// Handle failures here
	rc.Client.Set(key, jsonGame, 24*time.Hour)
	return nil
}

func GetGame(gameId string) (*models.Game, error) {

	key := GetGameKey(gameId)
	result, err := rc.Client.Get(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "Getting Game failed")
	}

	var game models.Game
	err = json.Unmarshal([]byte(result), &game)
	if err != nil {
		return nil, errors.Wrap(err, "Converting game meta json to game object failed")
	}

	return &game, nil
}

func SetGameMeta(gameMeta *models.GameMeta) error {

	key := GetGameMetaKey(gameMeta.GameId)

	jsonGame, err := json.Marshal(gameMeta)
	if err != nil {
		return errors.Wrap(err, "Game json conversion failed while setting game")
	}

	// Handle failures here
	rc.Client.Set(key, jsonGame, 24*time.Hour)
	return nil
}

func GetGameMeta(gameId string) (*models.GameMeta, error) {

	key := GetGameMetaKey(gameId)
	result, err := rc.Client.Get(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "Getting Game Meta failed")
	}

	var game models.GameMeta
	err = json.Unmarshal([]byte(result), &game)
	if err != nil {
		return nil, errors.Wrap(err, "Converting game meta json to game object failed")
	}

	return &game, nil
}

func SetGamePhrases(gameId string, newPhrases *models.PhraseList) error {
	key := fmt.Sprintf("game_phrases:%s", gameId) // Key pattern remains the same

	// Retrieve existing phrases
	existingPhrasesJSON, err := rc.Client.Get(key).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("Error getting existing phrases from Redis:", err)
		return err
	}

	var existingPhrases models.PhraseList
	if existingPhrasesJSON != "" {
		err = json.Unmarshal([]byte(existingPhrasesJSON), &existingPhrases)
		if err != nil {
			fmt.Println("Error unmarshaling existing phrases:", err)
			return err
		}
	} else {
		// Initialize existingPhrases.List if no phrases are currently stored
		existingPhrases.List = &[]models.Phrase{}
	}

	// Append new phrases to the existing phrases
	*existingPhrases.List = append(*existingPhrases.List, *newPhrases.List...)

	// Marshal the updated phrases list
	updatedPhrasesJSON, err := json.Marshal(existingPhrases)
	if err != nil {
		fmt.Println("Error marshaling updated phrases:", err)
		return err
	}

	// Save the updated list back to Redis
	err = rc.Client.Set(key, updatedPhrasesJSON, 24*time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting updated phrases in Redis:", err)
		return err
	}

	return nil
}

func GetGamePhrases(gameId string) (*models.PhraseList, error) {
	key := fmt.Sprintf("game_phrases:%s", gameId)

	result, err := rc.Client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("No phrases found for game", gameId)
			return nil, nil
		}
		fmt.Println("Error getting game phrases from Redis:", err)
		return nil, err
	}

	var phrases models.PhraseList
	err = json.Unmarshal([]byte(result), &phrases)
	if err != nil {
		fmt.Println("Error unmarshaling game phrases:", err)
		return nil, err
	}

	return &phrases, nil
}

func SetPlayerPhrases(playerId string, phrases *models.PhraseList) error {
	key := fmt.Sprintf("player_phrases:%s", playerId)
	jsonPhrases, err := json.Marshal(phrases)
	if err != nil {
		fmt.Println("Error marshaling phrases:", err)
		return err
	}

	err = rc.Client.Set(key, jsonPhrases, 24*time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting player game phrases in Redis:", err)
		return err
	}
	fmt.Println("Player Phrases set successfully")
	return nil
}

func GetPlayerPhrases(playerId string) (*models.PhraseList, error) {
	key := fmt.Sprintf("player_phrases:%s", playerId)

	result, err := rc.Client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("No phrases found for player in game", playerId)
			return nil, nil
		}
		fmt.Println("Error getting player game phrases from Redis:", err)
		return nil, err
	}

	var phrases models.PhraseList
	err = json.Unmarshal([]byte(result), &phrases)
	if err != nil {
		fmt.Println("Error unmarshaling player game phrases:", err)
		return nil, err
	}

	return &phrases, nil
}

func GetPlayerKey(playerId string) string {
	return fmt.Sprintf("player:%s", playerId)
}

func GetGameKey(gameId string) string {
	return fmt.Sprintf("game:%s", gameId)
}

func GetGameMetaKey(gameId string) string {
	return fmt.Sprintf("gamemeta:%s", gameId)
}

func GetUserInputKey(playerId string, gameId string) string {
	// Ideally we should use a different db like MySQL for storing words
	return fmt.Sprintf("phrases:%s:%s", gameId, playerId)
}
