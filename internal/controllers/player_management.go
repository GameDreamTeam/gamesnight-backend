package controllers

import (
	"errors"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetPlayerDetailsController(c *gin.Context) {
	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	player := p.(*models.Player)

	playerInfo, err := services.GetPlayerService().GetPlayerDetails(*player.Id)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, playerInfo, nil)
}

func GetPlayerPhrasesController(c *gin.Context) {
	playerId := c.Param("playerId")

	phrases, err := services.GetPlayerService().GetPlayerPhrases(playerId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	c.JSON(http.StatusOK, phrases)
}

func RemovePlayerController(c *gin.Context) {
	// Get player ID to be removed from the request
	playerId := c.Param("playerId")

	// Fetch the game meta using game ID
	gameId := c.Param("gameId")
	gameMeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	adminId := gameMeta.AdminId

	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	player := p.(*models.Player)

	if *player.Id != gameMeta.AdminId {
		logger.GetLogger().Logger.Error(
			"Only admin should remove players",
			zap.Any("gamemeta", gameMeta),
			zap.Any("player", player),
		)
		SendResponse(c, http.StatusInternalServerError, nil,
			errors.New("only admin should remove players"))
		return
	}

	// Validate that the player to be removed is not admin
	if adminId != playerId {

		// Remove the player from the game meta and write to redis
		updatedGameMeta, err := services.GetPlayerService().RemovePlayer(gameMeta, playerId)
		if err != nil {
			SendResponse(c, http.StatusInternalServerError, nil, err)
			return
		}

		SendResponse(c, http.StatusOK, updatedGameMeta, nil)
	} else {
		SendResponse(c, http.StatusBadRequest, nil, errors.New("bad Request: Admin cannot remove itself"))
	}
}
