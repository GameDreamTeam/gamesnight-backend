package database

import (
	"fmt"

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

func GetRedis() *RedisClient {
	return rc
}
