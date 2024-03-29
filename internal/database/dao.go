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
		return errors.New("game json conversion failed while setting game")
	}

	err = rc.Client.Set(key, jsonGame, 24*time.Hour).Err()
	if err != nil {
		return errors.Wrap(err, "Failed to set game meta in Redis")
	}
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

func GetPlayerDetails(playerId string) (*models.Player, error) {
	key := playerId
	result, err := rc.Client.Get(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "Getting Player failed")
	}

	var player models.Player
	err = json.Unmarshal([]byte(result), &player)
	if err != nil {
		return nil, errors.Wrap(err, "Converting player json to game object failed")
	}

	return &player, nil
}

func SetPlayerDetails(player models.Player) error {
	// Use player key for redis (keys.go)
	key := *player.Id
	jsonPlayer, err := json.Marshal(player)
	if err != nil {
		return errors.Wrap(err, "Player json conversion failed while setting game")
	}

	err = rc.Client.Set(key, jsonPlayer, 24*time.Hour).Err()
	if err != nil {
		return errors.Wrap(err, "Failed to set Player in Redis")
	}
	return nil
}

func SetGameMeta(gameMeta *models.GameMeta) error {
	key := GetGameMetaKey(gameMeta.GameId)

	jsonGame, err := json.Marshal(gameMeta)
	if err != nil {
		return errors.Wrap(err, "Game json conversion failed")
	}

	err = rc.Client.Set(key, jsonGame, 24*time.Hour).Err()
	if err != nil {
		return errors.Wrap(err, "Failed to set game in Redis")
	}
	return nil
}

func GetGameMeta(gameId string) (*models.GameMeta, error) {
	key := GetGameMetaKey(gameId)

	result, err := rc.Client.Get(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "game not found")
	}
	var gameMeta models.GameMeta
	err = json.Unmarshal([]byte(result), &gameMeta)
	if err != nil {
		return nil, errors.Wrap(err, "converting game meta json to game object failed")
	}

	return &gameMeta, nil
}

func SetGamePhrases(gameId string, newPhrases *models.PhraseList) error {
	key := GetGamePhraseKey(gameId)

	// At the very least check if same phrase is not entered twice, maybe using frontend?
	// Can add checks to whether submitted phrases exist already ( Find a way to do in realtime or upon submit )
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
		existingPhrases.List = &[]models.Phrase{}
	}

	*existingPhrases.List = append(*existingPhrases.List, *newPhrases.List...)

	updatedPhrasesJSON, err := json.Marshal(existingPhrases)
	if err != nil {
		fmt.Println("Error marshaling updated phrases:", err)
		return err
	}

	err = rc.Client.Set(key, updatedPhrasesJSON, 24*time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting updated phrases in Redis:", err)
		return err
	}

	return nil
}

func GetGamePhrases(gameId string) (*models.PhraseList, error) {
	key := GetGamePhraseKey(gameId)

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
	key := GetPlayerPhraseKey(playerId)
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
	key := GetPlayerPhraseKey(playerId)

	result, err := rc.Client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("No phrases found for player in game", playerId)
			return nil, err
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

func SetCurrentPhraseMap(gameId string, phraseStatusMap models.PhraseStatusMap) error {
	key := GetCurrentPhraseMapKey(gameId + "current")

	jsonMap, err := json.Marshal(phraseStatusMap)
	if err != nil {
		return errors.Wrap(err, "error marshaling PhraseStatusMap")
	}

	err = rc.Client.Set(key, jsonMap, 24*time.Hour).Err()
	if err != nil {
		return errors.Wrap(err, "failed to set PhraseStatusMap in Redis")
	}

	return nil
}

// This is business logic creeping into database access layer
func GetCurrentPhraseMap(gameId string) (models.PhraseStatusMap, error) {
	key := GetCurrentPhraseMapKey(gameId + "current")
	result, err := rc.Client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			// No PhraseStatusMap found for the game, return an empty map
			return models.PhraseStatusMap{}, nil
		}
		return models.PhraseStatusMap{}, errors.Wrap(err, "getting current phrase map failed")
	}

	var phraseStatusMap models.PhraseStatusMap
	err = json.Unmarshal([]byte(result), &phraseStatusMap)
	if err != nil {
		return models.PhraseStatusMap{}, errors.Wrap(err, "error unmarshaling PhraseStatusMap")
	}

	return phraseStatusMap, nil
}
