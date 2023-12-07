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
		SendResponse(c, http.StatusInternalServerError, nil,
			errors.New("player starting game should be admin"))
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

	PhraseMap, err := services.GetGameService().GeneratePhraseListToMap(gamePhrases)
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

func StartTurnController(c *gin.Context) {
	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)
	// Throw different error if game is not playing
	if err != nil || game.GameState != models.Playing {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	player := p.(*models.Player)
	models.CurrentIndex = 0

	if *player.Id != *game.CurrentPlayer.Id {
		logger.GetLogger().Logger.Error(
			"player starting turn should be current player",
			zap.Any("game", game),
			zap.Any("player", player),
		)
		SendResponse(c, http.StatusInternalServerError, nil,
			errors.New("player starting turn should be current player"))
		return
	}

	currentPhraseMap, err := services.GetGameService().GetCurrentPhraseMap(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	PhraseToBeGuessed, err := services.GetGameService().GetPhraseToBeGuessed(currentPhraseMap)

	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	responseData := models.ResponseData{
		PhraseMap:     &currentPhraseMap,
		CurrentPhrase: PhraseToBeGuessed,
	}

	services.GetGameService().StartTurnTimer(gameId)

	SendResponse(c, http.StatusOK, responseData, nil)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	player := p.(*models.Player)

	if *player.Id != gameMeta.AdminId {
		logger.GetLogger().Logger.Error(
			"Only admin should remove players",
			zap.Any("gamemeta", gameMeta),
			zap.Any("player", player),
		)
		SendResponse(c, http.StatusInternalServerError, nil,
			errors.New("only admin should remove players"))
		return
	}

	// Validate that the player to be removed is not admin
	if adminId != playerId {

		// Remove the player from the game meta and write to redis
		updatedGameMeta, err := services.GetGameService().RemovePlayer(gameMeta, playerId)
		if err != nil {
			SendResponse(c, http.StatusInternalServerError, nil, err)
			return
		}

		SendResponse(c, http.StatusOK, updatedGameMeta, nil)
	} else {
		SendResponse(c, http.StatusBadRequest, nil, errors.New("bad Request: Admin cannot remove itself"))
	}
}
