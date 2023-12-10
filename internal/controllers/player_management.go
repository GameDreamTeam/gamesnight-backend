package controllers

import (
	"errors"
	"gamesnight/internal/models"
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
	// Get player ID to be removed from the request
	playerId := c.Param("playerId")

	// Fetch the game meta using game ID
	gameId := c.Param("gameId")
	gameMeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	adminId := gameMeta.AdminId

	p, exists := c.Get("player")
	if !exists {
		SendResponse(c, http.StatusInternalServerError, nil, err)
	}
	player := p.(*models.Player)

	err = isAdminPlayer(*gameMeta, player)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
	}

	// Validate that the player to be removed is not admin
	if adminId != playerId {
		// Remove the player from the game meta and write to redis
		updatedGameMeta, err := services.GetPlayerService().RemovePlayer(gameMeta, playerId)
		if err != nil {
			SendResponse(c, http.StatusInternalServerError, nil, err)
			return
		}

		SendResponse(c, http.StatusOK, updatedGameMeta, nil)
	} else {
		SendResponse(c, http.StatusBadRequest, nil, errors.New("bad Request: Admin cannot remove itself"))
	}
}
