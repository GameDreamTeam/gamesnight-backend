package controllers

import (
	"errors"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayerDetailsController(c *gin.Context) {
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	playerInfo, err := services.GetPlayerService().GetPlayerDetails(*player.Id)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, playerInfo, nil)
}

func GetPlayerPhrasesController(c *gin.Context) {
	playerId := c.Param("playerId")

	phrases, err := services.GetPlayerService().GetPlayerPhrases(playerId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, phrases, nil)
}

func RemovePlayerController(c *gin.Context) {
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

	err = isAdminPlayer(*gameMeta, player)
	if err != nil {
		SendResponse(c, http.StatusForbidden, nil, err)
	}

	adminId := gameMeta.AdminId
	playerToBeRemovedId := c.Param("playerId")
	if adminId != playerToBeRemovedId {
		updatedGameMeta, err := services.GetPlayerService().RemovePlayer(gameMeta, playerToBeRemovedId)
		if err != nil {
			SendResponse(c, http.StatusNotFound, nil, err)
			return
		}

		SendResponse(c, http.StatusOK, updatedGameMeta, nil)
	} else {
		SendResponse(c, http.StatusBadRequest, nil, errors.New("bad Request: Admin cannot remove itself"))
	}
}
