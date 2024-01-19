package controllers

import (
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGameMetaController(c *gin.Context) {
	gameId := c.Param("gameId")
	gameMeta, err := services.GetGameService().GetGameMeta(gameId)

	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}
	SendResponse(c, http.StatusOK, gameMeta, nil)
}

func GetGameController(c *gin.Context) {
	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)

	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}
	SendResponse(c, http.StatusOK, game, nil)
}

func GetGamePhrasesController(c *gin.Context) {
	gameId := c.Param("gameId")

	phrases, err := services.GetGameService().GetGamePhrases(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, phrases, nil)
}

func GetCurrentGamePhrasesController(c *gin.Context) {
	gameId := c.Param("gameId")

	phrases, err := services.GetGameService().GetCurrentPhraseMap(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, phrases, nil)
}
