package controllers

import (
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartGameController(c *gin.Context) {
	gameId := c.Param("gameId")
	// Rename this to check if game exists or not
	// We need to always check if game exists or not. This can be a middleware infact
	gamemeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	err = isAdminPlayer(*gamemeta, player)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	game, err := services.GetGameService().StartGame(gamemeta.GameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	gamePhrases, err := services.GetGameService().GetGamePhrases(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	PhraseMap, err := services.GeneratePhraseListToMap(gamePhrases)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	err = services.GetGameService().SetCurrentPhraseMap(gameId, PhraseMap)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, game, nil)
}
