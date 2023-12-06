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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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

	c.JSON(http.StatusOK, gin.H{"success": "Phrases added successfully"})
}

func GetGamePhrasesController(c *gin.Context) {
	gameId := c.Param("gameId")

	phrases, err := services.GetGameService().GetGamePhrases(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	c.JSON(http.StatusOK, phrases)
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
