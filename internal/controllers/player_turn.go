package controllers

import (
	"errors"
	"fmt"
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
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

	err = services.GetGameService().CheckCurrentPlayer(*player.Id, *game.CurrentPlayer.Id)
	if err != nil {
		SendResponse(c, http.StatusForbidden, nil, errors.New("player starting turn should be current player"))
		return
	}

	currentPhraseMap, err := services.GetGameService().GetCurrentPhraseMap(*game)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	PhraseToBeGuessed, err := services.GetGameService().GetPhraseToBeGuessed(currentPhraseMap, *game)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, PhraseToBeGuessed, nil)
}

func PlayerGuessController(c *gin.Context) {
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)

	if err != nil || game.GameState != models.Playing {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	err = services.GetGameService().CheckCurrentPlayer(*player.Id, *game.CurrentPlayer.Id)
	if err != nil {
		SendResponse(c, http.StatusForbidden, nil, errors.New("player starting turn should be current player"))
		return
	}

	var guessRequest models.PlayerGuess
	err = BindJSONAndHandleError(c, &guessRequest)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	currentPhraseMap, err := services.GetGameService().HandlePlayerGuess(*game, guessRequest.PlayerChoice)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	game, err = services.GetGameService().GetGame(gameId)
	PhraseToBeGuessed, err := services.GetGameService().GetPhraseToBeGuessed(currentPhraseMap, *game)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, PhraseToBeGuessed, nil)
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

	err = services.GetGameService().CheckCurrentPlayer(*player.Id, *game.CurrentPlayer.Id)
	if err != nil {
		SendResponse(c, http.StatusForbidden, nil, errors.New("player ending turn should be current player"))
		return
	}

	_, err = services.GetPlayerService().NextPlayerAndTeam(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	currentPhraseMap, err := services.GetGameService().GetCurrentPhraseMap(*game)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	updatedPhraseMap := services.GetGameService().RemoveGuessedPhrases(gameId, currentPhraseMap)
	fmt.Println(updatedPhraseMap)
	SendResponse(c, http.StatusOK, updatedPhraseMap, nil)
}
