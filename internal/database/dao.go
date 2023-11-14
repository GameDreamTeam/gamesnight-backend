package database

import (
	"encoding/json"
	"fmt"
	"gamesnight/internal/models"
	"time"

	"github.com/pkg/errors"
)

func SetGame(game *models.Game) error {
	key := GetGameKey(game.GameId)

	jsonGame, err := json.Marshal(game)
	if err != nil {
		return errors.Wrap(err, "Game json conversion failed while setting game")
	}

	// Handle failures here
	rc.Client.Set(key, jsonGame, 24*time.Hour)
	return nil
}

func GetGame(gameId string) (*models.Game, error) {

	key := GetGameKey(gameId)
	result, err := rc.Client.Get(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "Getting Game Meta failed")
	}

	var game models.Game
	err = json.Unmarshal([]byte(result), &game)
	if err != nil {
		return nil, errors.Wrap(err, "Converting game meta json to game object failed")
	}

	return &game, nil
}

func SetGameMeta(gameMeta *models.GameMeta) error {

	key := GetGameMetaKey(gameMeta.GameId)

	jsonGame, err := json.Marshal(gameMeta)
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

func GetGameKey(playerId string) string {
	return fmt.Sprintf("game:%s", playerId)
}

func GetGameMetaKey(gameId string) string {
	return fmt.Sprintf("gamemeta:%s", gameId)
}

func GetUserInputKey(playerId string, gameId string) string {
	// Ideally we should use a different db like MySQL for storing words
	return fmt.Sprintf("phrases:%s:%s", gameId, playerId)
}
