package database

import (
	"fmt"
)

func GetCurrentPhraseMapKey(gameId string) string {
	return fmt.Sprintf("current-phrases:%s", gameId)
}

func GetGameKey(gameId string) string {
	return fmt.Sprintf("game:%s", gameId)
}

func GetGameMetaKey(gameId string) string {
	return fmt.Sprintf("gamemeta:%s", gameId)
}

func GetGamePhraseKey(gameId string) string {
	return fmt.Sprintf("game-phrase:%s", gameId)
}

func GetPlayerPhraseKey(playerId string) string {
	return fmt.Sprintf("player-phrase:%s", playerId)
}
