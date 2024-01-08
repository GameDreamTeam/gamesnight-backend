package controllers

import (
	"errors"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StartTurnController(c *gin.Context) {
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	if game.GameState != models.Playing {
		SendResponse(c, http.StatusConflict, nil, errors.New("the game is not in the playing State"))
		return
	}

	// I think this could also be a middleware
	// Also note middleware is also like abstracting away some logic into a function
	// It makes it clear what requirements are needed

	err = services.GetGameService().CheckCurrentPlayer(*player.Id, game.GameId)
	if err != nil {
		SendResponse(c, http.StatusForbidden, nil, errors.New("player starting turn should be current player"))
		return
	}

	currentPhraseMap, err := services.GetGameService().GetCurrentPhraseMap(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	PhraseToBeGuessed, err := services.GetGameService().GetPhraseToBeGuessed(currentPhraseMap, game.CurrentPhraseMapIndex)

	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	responseData := models.ResponseData{
		PhraseMap:     &currentPhraseMap,
		CurrentPhrase: PhraseToBeGuessed,
		TurnStartedAt: models.TurnStartTime,
	}

	SendResponse(c, http.StatusOK, responseData, nil)
}

func EndTurnController(c *gin.Context) {
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	if game.GameState != models.Playing {
		SendResponse(c, http.StatusConflict, nil, errors.New("the game is not in the playing State"))
		return
	}

	err = services.GetGameService().CheckCurrentPlayer(*player.Id, game.GameId)
	if err != nil {
		SendResponse(c, http.StatusForbidden, nil, errors.New("player ending turn should be current player"))
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

	updatedPhraseMap := services.GetGameService().RemoveGuessedPhrases(gameId, currentPhraseMap)

	models.TurnStartTime = time.Time{}
	SendResponse(c, http.StatusOK, updatedPhraseMap, nil)
}

func PlayerGuessController(c *gin.Context) {
	// can remove validations here to make API lightweight
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)

	// Throw different error if game is not playing
	if err != nil || game.GameState != models.Playing {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

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
	var guessRequest models.PlayerGuess
	err = BindJSONAndHandleError(c, &guessRequest)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
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

	PhraseToBeGuessed, err := services.GetGameService().GetPhraseToBeGuessed(currentPhraseMap, game.CurrentPhraseMapIndex)

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
