package controllers

import (
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeTeamsController(c *gin.Context) {
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

	err = services.GetGameService().CheckIfAllPlayerHaveSubmittedPhrases(*gameMeta)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	err = isAdminPlayer(*gameMeta, player)
	if err != nil {
		SendResponse(c, http.StatusForbidden, nil, err)
		return
	}

	gameWithTeams, err := services.GetGameService().MakeTeams(gameMeta)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, gameWithTeams, nil)
}
