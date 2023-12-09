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

func StartTurnController(c *gin.Context) {
	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)
	// Throw different error if game is not playing
	if err != nil || game.GameState != models.Playing {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	player := p.(*models.Player)
	models.CurrentIndex = 0

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

	currentPhraseMap, err := services.GetGameService().GetCurrentPhraseMap(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	PhraseToBeGuessed, err := services.GetGameService().GetPhraseToBeGuessed(currentPhraseMap)

	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	responseData := models.ResponseData{
		PhraseMap:     &currentPhraseMap,
		CurrentPhrase: PhraseToBeGuessed,
	}

	// services.GetGameService().StartTurnTimer(gameId)

	SendResponse(c, http.StatusOK, responseData, nil)
}

func EndTurnController(c *gin.Context) {
	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)
	// Throw different error if game is not playing
	if err != nil || game.GameState != models.Playing {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	player := p.(*models.Player)
	if *player.Id != *game.CurrentPlayer.Id {
		logger.GetLogger().Logger.Error(
			"player ending turn should be current player",
			zap.Any("game", game),
			zap.Any("player", player),
		)
		SendResponse(c, http.StatusInternalServerError, nil,
			errors.New("player ending turn should be current player"))
		return
	}

	_, err = services.GetPlayerService().NextPlayerAndTeam(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}
	currentPhraseMap, err := services.GetGameService().GetCurrentPhraseMap(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}
	SendResponse(c, http.StatusOK, currentPhraseMap, nil)
}

func PlayerGuessController(c *gin.Context) {
	// can remove validations here to make API lightweight
	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)

	// Throw different error if game is not playing
	if err != nil || game.GameState != models.Playing {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	player := p.(*models.Player)

	if *player.Id != *game.CurrentPlayer.Id {
		logger.GetLogger().Logger.Error(
			"player making guess should be current player",
			zap.Any("game", game),
			zap.Any("player", player),
		)
		SendResponse(c, http.StatusInternalServerError, nil,
			errors.New("player making guess should be current player"))
		return
	}

	// Parse request body
	var guessRequest models.PlayerGuessWithWord
	if err := c.BindJSON(&guessRequest); err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	err = services.GetGameService().HandlePlayerGuess(*game, guessRequest.PlayerChoice)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	currentPhraseMap, err := services.GetGameService().GetCurrentPhraseMap(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	PhraseToBeGuessed, err := services.GetGameService().GetPhraseToBeGuessed(currentPhraseMap)

	if err != nil {
		// SendResponse(c, http.StatusInternalServerError, nil, err)
		//Write End Game Service here
		return
	}

	responseData := models.ResponseData{
		PhraseMap:     &currentPhraseMap,
		CurrentPhrase: PhraseToBeGuessed,
	}

	SendResponse(c, http.StatusOK, responseData, nil)

}
