package database

import (
	"encoding/json"
	"fmt"
	"gamesnight/internal/models"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

type RedisClient struct {
	Client *redis.Client
}

var rc *RedisClient

func NewRedisClient() {
	mr, err := miniredis.Run()
	if err != nil {
		fmt.Println(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	rc = &RedisClient{
		Client: rdb,
	}
}

func SetGame(game *models.Game) error {

	key := GetGameKey(game.GameId)

	jsonGame, err := json.Marshal(game)
	if err != nil {
		fmt.Println("Error marshaling game:", err)
		return err
	}

	rc.Client.Set(key, jsonGame, 24*time.Hour)
	return nil
}

func GetGame(gameId string) (*models.Game, error) {

	key := GetGameKey(gameId)
	result, err := rc.Client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	var game models.Game
	err = json.Unmarshal([]byte(result), &game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func SetGamePhrases(gameId string, phrases *models.PhraseList) error {
	key := fmt.Sprintf("game_phrases:%s", gameId) // Key pattern can be "game_phrases:<gameId>"

	jsonPhrases, err := json.Marshal(phrases)
	if err != nil {
		fmt.Println("Error marshaling phrases:", err)
		return err
	}

	err = rc.Client.Set(key, jsonPhrases, 24*time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting game phrases in Redis:", err)
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

func SetPlayerGamePhrases(gameId string, playerId string, phrases *models.PhraseList) error {
	key := fmt.Sprintf("player_game_phrases:%s:%s", gameId, playerId) // Key pattern can be "player_game_phrases:<gameId>:<playerId>"

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

	return nil
}

func GetPlayerGamePhrases(gameId string, playerId string) (*models.PhraseList, error) {
	key := fmt.Sprintf("player_game_phrases:%s:%s", gameId, playerId)

	result, err := rc.Client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("No phrases found for player in game", gameId, playerId)
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
