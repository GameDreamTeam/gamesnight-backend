package database

import (
	"fmt"
	"gamesnight/internal/config"

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

	if config.Get().Env == "local" {
		rdb, err := getLocalClient()

		if err != nil {
			panic(fmt.Sprintf("Redis Initialization failed for env%s with error %s", config.Get().Env, err))
		}

		rc = &RedisClient{
			Client: rdb,
		}
	}
}
