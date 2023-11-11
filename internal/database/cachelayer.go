package database

import (
	"fmt"
)

type CacheLayer struct{}

var cl *CacheLayer

func NewCacheLayer() {
	cl = &CacheLayer{}
}

func GetCacheLayer() *CacheLayer {
	return cl
}

func GetUserKey(userId string) string {
	return fmt.Sprintf("user:%s", userId)
}

func GetGameKey(gameId string) string {
	return fmt.Sprintf("game:%s", gameId)
}
