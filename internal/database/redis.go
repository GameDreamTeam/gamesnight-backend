package database

import (
	"encoding/json"
	"fmt"
	"gamesnight/internal/models"
	"time"

	"github.com/pkg/errors"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

type RedisClient struct {
	Client *redis.Client
}

var rc *RedisClient

func getLocalClient() (*redis.Client, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return rdb, nil
}

func NewRedisClient() {

	rdb, err := getLocalClient()

	if err != nil {
		panic(fmt.Sprintf("Redis Initialization failed %s", err))
	}

	rc = &RedisClient{
		Client: rdb,
	}
}

func SetGameMeta(game *models.GameMeta) error {

	key := GetGameMetaKey(game.GameId)

	jsonGame, err := json.Marshal(game)
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

func GetPlayerKey(playerId string) string {
	return fmt.Sprintf("player:%s", playerId)
}

func GetGameMetaKey(gameId string) string {
	return fmt.Sprintf("gamemeta:%s", gameId)
}

func GetUserInputKey(playerId string, gameId string) string {
	// Ideally we should use a different db like MySQL for storing words
	return fmt.Sprintf("phrases:%s:%s", gameId, playerId)
}
