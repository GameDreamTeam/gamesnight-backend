package controllers

import (
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGameMetaController(c *gin.Context) {
	gameId := c.Param("gameId")
	gameMeta, err := services.GetGameService().GetGameMeta(gameId)

	// If we don't find a game meta with a gameId then we should return 404
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}
	SendResponse(c, http.StatusOK, gameMeta, nil)
}

func GetGameController(c *gin.Context) {
	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)

	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}
	SendResponse(c, http.StatusOK, game, nil)
}
