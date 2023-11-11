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
