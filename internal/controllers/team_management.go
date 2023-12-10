package controllers

import (
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeTeamsController(c *gin.Context) {
	//Check if player exist
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	//Check if game exist
	gameId := c.Param("gameId")
	gamemeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	//Check if make teams is called by admin
	err = isAdminPlayer(*gamemeta, player)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
	}

	game, err := services.GetGameService().MakeTeams(gamemeta)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, game, nil)
}
