package database

import (
	"fmt"
	"gamesnight/internal/models"
)

type CacheLayer struct{}

var cl *CacheLayer

func NewCacheLayer() {
	cl = &CacheLayer{}
}

func GetCacheLayer() *CacheLayer {
	return cl
}

func GetUserKey(user models.User) string {
	return fmt.Sprintf("user:%s", user.UserName)
}

func GetGameKey(game models.Game) string {
	return fmt.Sprintf("game:%s", game.GameId)
}
