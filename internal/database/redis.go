package database

import (
	"fmt"
	"gamesnight/internal/config"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"os"
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
	var rdb *redis.Client
    var err error

	if config.Get().Env == "local" {
		rdb, err = getLocalClient()

	}else{
		redisAddr := os.Getenv("REDIS_ADDR") // e.g., "prod-redis-host:6379"
        redisPassword := os.Getenv("REDIS_PASSWORD") // e.g., "mysecretpassword"

        rdb = redis.NewClient(&redis.Options{
            Addr:     redisAddr,
            Password: redisPassword, // no password set for empty string
            DB:       0, // default DB
        })
	}
	
	if err != nil {
		panic(fmt.Sprintf("Redis Initialization failed for env%s with error %s", config.Get().Env, err))
	}

	rc = &RedisClient{
		Client: rdb,
	}
}
