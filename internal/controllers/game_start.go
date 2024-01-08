package controllers

import (
	"errors"
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartGameController(c *gin.Context) {
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	gameId := c.Param("gameId")
	gamemeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	err = isAdminPlayer(*gamemeta, player)
	if err != nil {
		SendResponse(c, http.StatusForbidden, nil, err)
		return
	}

	game, err := services.GetGameService().GetGame(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	if game.GameState != models.TeamsDivided {
		SendResponse(c, http.StatusConflict, nil, errors.New("The game has not been divided into teams"))
		return
	}

	gameWithInitialisation, err := services.GetGameService().StartGame(game)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	gamePhrases, err := services.GetGameService().GetGamePhrases(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
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

	SendResponse(c, http.StatusOK, gameWithInitialisation, nil)
}
