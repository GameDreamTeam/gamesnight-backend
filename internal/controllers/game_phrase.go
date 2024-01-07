package controllers

import (
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPhraseController(c *gin.Context) {
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	gameId := c.Param("gameId")

	err = services.GetPlayerService().PlayerExistInGame(gameId, *player)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	err = services.GetPlayerService().PlayerAlreadyAddedPhrases(*player.Id)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	var phraseList models.PhraseList
	err = BindJSONAndHandleError(c, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	err = CheckPhraseListLength(phraseList)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	err = services.GetGameService().AddPhrasesToGame(*player.Id, gameId, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	err = services.GetGameService().AddPhrasesToPlayer(*player, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, phraseList, nil)
}

func GetGamePhrasesController(c *gin.Context) {
	gameId := c.Param("gameId")

	phrases, err := services.GetGameService().GetGamePhrases(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, phrases, nil)
}

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
