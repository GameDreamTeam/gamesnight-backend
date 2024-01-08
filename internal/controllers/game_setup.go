package controllers

import (
	"errors"
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewGameController(c *gin.Context) {
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	gameMeta, err := services.GetGameService().CreateNewGame(*player.Id)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, gameMeta, nil)
}

func JoinGameController(c *gin.Context) {
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	var username models.PlayerName
	err = BindJSONAndHandleError(c, &username)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}
	player.Name = &username.Name
	player.PhrasesSubmitted = false

	gameId := c.Param("gameId")

	game, err := services.GetGameService().GetGame(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}
	if game.GameState != models.PlayersJoining {
		SendResponse(c, http.StatusConflict, nil, errors.New("The game is not in the joiningState"))
		return
	}

	gameMeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	gameMetaWithPlayer, err := services.GetGameService().JoinGame(gameMeta, player)
	if err != nil {
		SendResponse(c, http.StatusConflict, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, gameMetaWithPlayer, nil)
}

func UpdateState(c *gin.Context) {
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	gameId := c.Param("gameId")
	gameMeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	if len(*gameMeta.Players) < 2 {
		SendResponse(c, http.StatusBadRequest, nil, errors.New("minimum 2 players to continue with game"))
		return
	}

	err = isAdminPlayer(*gameMeta, player)
	if err != nil {
		SendResponse(c, http.StatusForbidden, nil, err)
		return
	}

	game, err := services.GetGameService().UpdateStateOfGame(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, game, nil)
}
