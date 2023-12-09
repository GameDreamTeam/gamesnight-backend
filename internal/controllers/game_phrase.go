package controllers

import (
	"errors"
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPhraseController(c *gin.Context) {
	p, exists := c.Get("player")
	if !exists {
		SendResponse(c, http.StatusInternalServerError, nil, errors.New("internal server error"))
		return
	}

	var phraseList models.PhraseList
	player := p.(*models.Player)
	playerId := *player.Id
	gameId := c.Param("gameId")

	if err := c.BindJSON(&phraseList); err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	if len(*phraseList.List) != 4 {
		SendResponse(c, http.StatusBadRequest, nil, errors.New("total length of phrases must be 4"))
		return
	}

	err := services.GetGameService().AddPhrasesToGame(gameId, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	err = services.GetGameService().AddPhrasesToPlayer(playerId, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK,phraseList,nil)
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
