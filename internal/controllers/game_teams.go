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

func MakeTeamsController(c *gin.Context) {
	gameId := c.Param("gameId")
	gamemeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	player := p.(*models.Player)

	if *player.Id != gamemeta.AdminId {
		logger.GetLogger().Logger.Error(
			"player starting/making game/team should be admin",
			zap.Any("gamemeta", gamemeta),
			zap.Any("player", player),
		)
		SendResponse(c, http.StatusInternalServerError, nil,
			errors.New("player starting game should be admin"))
		return
	}

	game, err := services.GetGameService().MakeTeams(gamemeta)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, game, nil)
}
