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

func GeneratePhraseController(c *gin.Context) {
	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	player := p.(*models.Player)

	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)
	// Throw different error if game is not playing
	if err != nil || game.GameState != models.Playing {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	if *player.Id != *game.CurrentPlayer.Id {
		logger.GetLogger().Logger.Error(
			"player starting turn should be current player",
			zap.Any("game", game),
			zap.Any("player", player),
		)
		SendResponse(c, http.StatusInternalServerError, nil,
			errors.New("player starting turn should be current player"))
		return
	}

	// Get a random phrase for the player's turn
	phrases, err := services.GetGameService().GetRandomizedGamePhrases(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, phrases, nil)
}
