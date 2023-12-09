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

func NewGameController(c *gin.Context) {
	p, exists := c.Get("player")
	if !exists {
		SendResponse(c, http.StatusInternalServerError, nil, errors.New("internal Server Error"))
		return
	}

	//(*models.Player) is the type to which you are asserting that p should be converted.
	player := p.(*models.Player)

	game, err := services.GetGameService().CreateNewGame(*player.Id)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}
	SendResponse(c, http.StatusOK, game, nil)
}

func JoinGameController(c *gin.Context) {
	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	player := p.(*models.Player)

	gameId := c.Param("gameId")

	var playerName models.PlayerName

	if err := c.BindJSON(&playerName); err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	player.Name = &playerName.Username

	game, err := services.GetGameService().JoinGame(gameId, player)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}
	SendResponse(c, http.StatusOK, game, nil)
}

func StartGameController(c *gin.Context) {
	//Most of this below code should be part of some admin check middleware
	gameId := c.Param("gameId")
	gamemeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	player := p.(*models.Player)

	if *player.Id != gamemeta.AdminId {
		logger.GetLogger().Logger.Error(
			"player starting game should be admin",
			zap.Any("gamemeta", gamemeta),
			zap.Any("player", player),
		)
		SendResponse(c, http.StatusInternalServerError, nil, errors.New("player starting game should be admin"))
		return
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

	// Write PhraseMap to Redis
	err = services.GetGameService().SetCurrentPhraseMap(gameId, PhraseMap)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, game, nil)
}
