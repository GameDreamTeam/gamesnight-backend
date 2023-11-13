package controllers

import (
	"gamesnight/internal/database"
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

	// Can check if this type conversion is passing or failing
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

func GetGameMetaController(c *gin.Context) {
	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGameMeta(gameId)

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

	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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

	game, err := services.GetGameService().StartGame(gamemeta)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
	}

	SendResponse(c, http.StatusOK, game, nil)
}

func AddPhraseController(c *gin.Context) {
	var phraseList models.PhraseList
	gameId := c.Param("gameId")

	// Assuming the player ID is available in the context
	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	player := p.(*models.Player)
	playerId := *player.Id

	// Bind the incoming JSON to phraseList
	if err := c.BindJSON(&phraseList); err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	// Save the phrases in Redis with gameId as the key
	err := database.SetGamePhrases(gameId, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	// Optionally, save the phrases with a key combining the gameId and playerId
	err = database.SetPlayerGamePhrases(gameId, playerId, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Phrases added successfully"})
}

func GetGamePhrasesController(c *gin.Context) {
	gameId := c.Param("gameId")
	phrases, err := database.GetGamePhrases(gameId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	c.JSON(http.StatusOK, phrases)
}

func GetPlayerPhrasesController(c *gin.Context) {
	playerId := c.Param("playerId")
	gameId := c.Query("gameId") // Assuming gameId is also required to fetch specific phrases

	phrases, err := database.GetPlayerGamePhrases(gameId, playerId)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	c.JSON(http.StatusOK, phrases)
}
