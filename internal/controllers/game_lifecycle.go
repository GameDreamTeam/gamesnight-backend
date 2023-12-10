package controllers

import (
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewGameController(c *gin.Context) {
	//Check if player exist
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	//Create new game
	gameMeta, err := services.GetGameService().CreateNewGame(*player.Id)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, gameMeta, nil)
}

func JoinGameController(c *gin.Context) {
	//Check if player exist
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	var username models.PlayerName
	err = BindJSONAndHandleError(c, &username)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	player.Name = &username.Name

	gameId := c.Param("gameId")
	game, err := services.GetGameService().JoinGame(gameId, player)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, game, nil)
}

func StartGameController(c *gin.Context) {
	gameId := c.Param("gameId")
	gamemeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	err = isAdminPlayer(*gamemeta, player)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
	}

	game, err := services.GetGameService().StartGame(gamemeta.GameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	gamePhrases, err := services.GetGameService().GetGamePhrases(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	PhraseMap, err := services.GeneratePhraseListToMap(gamePhrases)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	err = services.GetGameService().SetCurrentPhraseMap(gameId, PhraseMap)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, game, nil)
}
