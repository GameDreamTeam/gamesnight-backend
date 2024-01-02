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
	// This should be part of check admin middleware
	err = isAdminPlayer(*gamemeta, player)
	if err != nil {
		// Error should be 403 instead of 500
		// We recognize the player and we feel they are not authorized to start the game
		SendResponse(c, http.StatusInternalServerError, nil, err)
	}

	game, err := services.GetGameService().MakeTeams(gamemeta)
	// So what we should do is, we should extend the basic error class provided by golang
	// and make custom errors and we could have a global error handler which figures out the
	// status code using the error type
	// The above call could return error because of multiple reasons other than just server error
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, game, nil)
}
