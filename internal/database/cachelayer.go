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

func GetPlayerKey(playerId string) string {
	return fmt.Sprintf("user:%s", playerId)
}

func GetGameKey(gameId string) string {
	return fmt.Sprintf("game:%s", gameId)
}

func GetUserInputKey(playerId string, gameId string) string {
	// Ideally we should use a different db like MySQL for storing words
	return fmt.Sprintf("input:%s:%s", gameId, playerId)
}
