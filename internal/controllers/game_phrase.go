package controllers

import (
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPhraseController(c *gin.Context) {
	//Check of Player exist
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	var phraseList models.PhraseList
	//Take phrases as an input from user
	err = BindJSONAndHandleError(c, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	//Check total number of phrases submitted
	err = CheckPhraseListLength(phraseList)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	gameId := c.Param("gameId")
	err = services.GetGameService().AddPhrasesToGame(gameId, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	playerId := *player.Id
	err = services.GetGameService().AddPhrasesToPlayer(playerId, &phraseList)
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
